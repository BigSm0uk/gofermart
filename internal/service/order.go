package service

import (
	"context"
	"fmt"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/google/uuid"
)

// OrderService представляет сервис для работы с заказами
type OrderService struct {
	orderRepo *repo.OrderRepository
}

// NewOrderService создает новый сервис заказов
func NewOrderService(orderRepo *repo.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

// CreateOrder создает новый заказ или возвращает существующий
// Возвращает order, isNew (true если заказ был создан, false если уже существовал), error
func (s *OrderService) CreateOrder(ctx context.Context, orderNumber string, userID uuid.UUID) (*domain.Order, bool, error) {
	// Пытаемся создать заказ
	order, isNew, err := s.orderRepo.CreateOrder(ctx, orderNumber, userID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create order: %w", err)
	}

	// Если заказ уже существовал, проверим владельца
	if !isNew && order.UserID != userID {
		// Заказ принадлежит другому пользователю
		return nil, false, domain.ErrOrderBelongsToOtherUser
	}

	return order, isNew, nil
}

// GetUserOrders получает все заказы пользователя
func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]domain.OrderResponse, error) {
	orders, err := s.orderRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Преобразуем в response формат
	responses := make([]domain.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = domain.OrderResponse{
			Number:     order.Number,
			Status:     order.Status,
			UploadedAt: order.UploadedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339
		}

		// Добавляем accrual только если он больше 0
		if order.Accrual > 0 {
			responses[i].Accrual = &order.Accrual
		}
	}

	return responses, nil
}

// UpdateOrderStatus обновляет статус заказа (используется воркером)
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderNumber string, status domain.OrderStatus, accrual *float64) error {
	return s.orderRepo.UpdateOrderStatus(ctx, orderNumber, status, accrual)
}

// GetOrdersForProcessing получает заказы для обработки воркером
func (s *OrderService) GetOrdersForProcessing(ctx context.Context, limit int) ([]domain.Order, error) {
	return s.orderRepo.GetOrdersForProcessing(ctx, limit)
}
