package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BigSm0uk/gofermart/internal/domain"
)

// OrderProcessor представляет процессор заказов
type OrderProcessor struct {
	orderService   *OrderService
	loyaltyService *LoyaltyService
	interval       time.Duration
	accrualURL     string

	// Rate limiting состояние
	rateLimitUntil time.Time    // время до которого приостановлены запросы
	rateLimitMutex sync.RWMutex // защита от race condition при чтении/записи rateLimitUntil

	// Счетчик активных воркеров для graceful shutdown
	activeWorkers int64
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
			// Проверяем, не приостановлены ли запросы из-за rate limiting
			if p.isRateLimited() {
				log.Printf("Rate limited, skipping processing cycle")
				continue
			}
			p.processOrders(ctx)
		}
	}
}

// isRateLimited проверяет, приостановлены ли запросы из-за rate limiting
func (p *OrderProcessor) isRateLimited() bool {
	p.rateLimitMutex.RLock()
	defer p.rateLimitMutex.RUnlock()
	return time.Now().Before(p.rateLimitUntil)
}

// setRateLimit устанавливает период приостановки запросов
func (p *OrderProcessor) setRateLimit(duration time.Duration) {
	p.rateLimitMutex.Lock()
	defer p.rateLimitMutex.Unlock()
	p.rateLimitUntil = time.Now().Add(duration)
	log.Printf("Rate limit set for %v until %v", duration, p.rateLimitUntil.Format(time.RFC3339))
}

// incrementActiveWorkers увеличивает счетчик активных воркеров
func (p *OrderProcessor) incrementActiveWorkers() {
	atomic.AddInt64(&p.activeWorkers, 1)
}

// decrementActiveWorkers уменьшает счетчик активных воркеров
func (p *OrderProcessor) decrementActiveWorkers() {
	atomic.AddInt64(&p.activeWorkers, -1)
}

// WaitForActiveWorkers ждет завершения всех активных воркеров
func (p *OrderProcessor) WaitForActiveWorkers(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&p.activeWorkers) == 0 {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// parseRetryAfter парсит заголовок Retry-After и возвращает длительность
func parseRetryAfter(retryAfter string) time.Duration {
	if retryAfter == "" {
		return 60 * time.Second // По умолчанию 60 секунд
	}

	// Пытаемся парсить как число секунд
	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Пытаемся парсить как HTTP дату (RFC1123)
	if retryTime, err := time.Parse(time.RFC1123, retryAfter); err == nil {
		duration := time.Until(retryTime)
		if duration > 0 {
			return duration
		}
	}

	// Если не удалось распарсить, возвращаем дефолтное значение
	log.Printf("Failed to parse Retry-After header: %s, using default 60s", retryAfter)
	return 60 * time.Second
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
		// Проверяем контекст перед каждой обработкой
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping order processing")
			return
		default:
		}

		// Проверяем rate limiting перед каждым заказом
		if p.isRateLimited() {
			log.Printf("Rate limited during processing, stopping")
			return
		}

		p.incrementActiveWorkers()
		p.processOrder(ctx, order)
		p.decrementActiveWorkers()
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

	// Обработка 429 Too Many Requests
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		duration := parseRetryAfter(retryAfter)
		p.setRateLimit(duration)
		log.Printf("Rate limited (429) from accrual API, setting rate limit for %v", duration)
		return fmt.Errorf("rate limited: %w", domain.ErrRateLimited)
	}

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

	// Обработка 429 Too Many Requests
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		duration := parseRetryAfter(retryAfter)
		p.setRateLimit(duration)
		log.Printf("Rate limited (429) from accrual API during status check, setting rate limit for %v", duration)
		return 0, domain.OrderStatusProcessing // Оставляем в обработке
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
