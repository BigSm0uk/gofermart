package service

import (
	"context"
	"log"
	"math/rand"
	"strconv"
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

	// Если заказ NEW, меняем статус на PROCESSING
	if order.Status == domain.OrderStatusNew {
		err := p.orderService.UpdateOrderStatus(ctx, order.Number, domain.OrderStatusProcessing, nil)
		if err != nil {
			log.Printf("Failed to update order %s status to PROCESSING: %v", order.Number, err)
			return
		}
		log.Printf("Order %s status updated to PROCESSING", order.Number)
	}

	// Эмулируем обращение к системе начисления баллов
	accrual, status := p.emulateAccrualSystem(order.Number)

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

// emulateAccrualSystem эмулирует работу системы начисления баллов
func (p *OrderProcessor) emulateAccrualSystem(orderNumber string) (float64, domain.OrderStatus) {
	// 10% заказов помечаем как невалидные
	if rand.Float64() < 0.1 {
		return 0, domain.OrderStatusInvalid
	}

	// Вычисляем начисление на основе номера заказа
	// Используем простую формулу: (сумма цифр % 10000) / 100
	sum := 0
	for _, char := range orderNumber {
		if char >= '0' && char <= '9' {
			digit, _ := strconv.Atoi(string(char))
			sum += digit
		}
	}

	accrual := float64(sum%10000) / 100.0

	return accrual, domain.OrderStatusProcessed
}
