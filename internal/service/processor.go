package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BigSm0uk/gofermart/internal/domain"
)

// OrderProcessor представляет процессор заказов
type OrderProcessor struct {
	orderService   *OrderService
	loyaltyService *LoyaltyService
	interval       time.Duration
	accrualURL     string
}

// NewOrderProcessor создает новый процессор заказов
func NewOrderProcessor(orderService *OrderService, loyaltyService *LoyaltyService, interval time.Duration, accrualURL string) *OrderProcessor {
	return &OrderProcessor{
		orderService:   orderService,
		loyaltyService: loyaltyService,
		interval:       interval,
		accrualURL:     accrualURL,
	}
}

// Run запускает процессор заказов
func (p *OrderProcessor) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	log.Printf("Order processor started with interval %v", p.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Order processor stopped")
			return
		case <-ticker.C:
			p.processOrders(ctx)
		}
	}
}

// processOrders обрабатывает заказы со статусом NEW и PROCESSING
func (p *OrderProcessor) processOrders(ctx context.Context) {
	// Получаем заказы для обработки
	orders, err := p.orderService.GetOrdersForProcessing(ctx, 100)
	if err != nil {
		log.Printf("Failed to get orders for processing: %v", err)
		return
	}

	if len(orders) == 0 {
		return
	}

	log.Printf("Processing %d orders", len(orders))

	for _, order := range orders {
		p.processOrder(ctx, order)
	}
}

// processOrder обрабатывает один заказ
func (p *OrderProcessor) processOrder(ctx context.Context, order domain.Order) {
	log.Printf("Processing order %s (status: %s)", order.Number, order.Status)

	// Если заказ NEW, регистрируем его в accrual системе
	if order.Status == domain.OrderStatusNew {
		err := p.registerOrderInAccrual(ctx, order.Number)
		if err != nil {
			log.Printf("Failed to register order %s in accrual system: %v", order.Number, err)
			return
		}

		err = p.orderService.UpdateOrderStatus(ctx, order.Number, domain.OrderStatusProcessing, nil)
		if err != nil {
			log.Printf("Failed to update order %s status to PROCESSING: %v", order.Number, err)
			return
		}
		log.Printf("Order %s registered in accrual system and status updated to PROCESSING", order.Number)
	}

	// Обращаемся к системе начисления баллов
	accrual, status := p.getAccrualFromAPI(ctx, order.Number)

	// Обновляем статус заказа
	err := p.orderService.UpdateOrderStatus(ctx, order.Number, status, &accrual)
	if err != nil {
		log.Printf("Failed to update order %s status: %v", order.Number, err)
		return
	}

	// Если заказ обработан успешно, начисляем баллы
	if status == domain.OrderStatusProcessed && accrual > 0 {
		err = p.loyaltyService.ProcessOrderAccrual(ctx, order.UserID, order.Number, accrual)
		if err != nil {
			log.Printf("Failed to process accrual for order %s: %v", order.Number, err)
		} else {
			log.Printf("Processed accrual %.2f for order %s", accrual, order.Number)
		}
	}

	log.Printf("Order %s processed with status %s", order.Number, status)
}

// registerOrderInAccrual регистрирует заказ в accrual системе
func (p *OrderProcessor) registerOrderInAccrual(ctx context.Context, orderNumber string) error {
	url := fmt.Sprintf("%s/api/orders", p.accrualURL)

	requestBody := map[string]string{
		"order": orderNumber,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to register order: %w", err)
	}
	defer resp.Body.Close()

	// 202 - заказ принят, 409 - заказ уже существует (это нормально)
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusConflict {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("accrual API returned status %d: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusConflict {
		log.Printf("Order %s already exists in accrual system", orderNumber)
	} else {
		log.Printf("Order %s registered in accrual system", orderNumber)
	}
	return nil
}

// AccrualResponse представляет ответ от accrual API
type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

// getAccrualFromAPI получает информацию о заказе от accrual API
func (p *OrderProcessor) getAccrualFromAPI(ctx context.Context, orderNumber string) (float64, domain.OrderStatus) {
	url := fmt.Sprintf("%s/api/orders/%s", p.accrualURL, orderNumber)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request for order %s: %v", orderNumber, err)
		return 0, domain.OrderStatusProcessing // Оставляем в обработке при ошибке
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get accrual info for order %s: %v", orderNumber, err)
		return 0, domain.OrderStatusProcessing // Оставляем в обработке при ошибке
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response for order %s: %v", orderNumber, err)
		return 0, domain.OrderStatusProcessing
	}

	// Если заказ не найден (404), возвращаем null - это нормально для новых заказов
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("Order %s not found in accrual system", orderNumber)
		return 0, domain.OrderStatusProcessing
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Accrual API returned status %d for order %s", resp.StatusCode, orderNumber)
		return 0, domain.OrderStatusProcessing
	}

	var accrualResp AccrualResponse
	if err := json.Unmarshal(body, &accrualResp); err != nil {
		log.Printf("Failed to parse accrual response for order %s: %v", orderNumber, err)
		return 0, domain.OrderStatusProcessing
	}

	// Конвертируем статус из accrual API в наш статус
	switch accrualResp.Status {
	case "REGISTERED":
		return 0, domain.OrderStatusProcessing
	case "PROCESSING":
		return 0, domain.OrderStatusProcessing
	case "INVALID":
		return 0, domain.OrderStatusInvalid
	case "PROCESSED":
		return accrualResp.Accrual, domain.OrderStatusProcessed
	default:
		log.Printf("Unknown status from accrual API: %s for order %s", accrualResp.Status, orderNumber)
		return 0, domain.OrderStatusProcessing
	}
}
