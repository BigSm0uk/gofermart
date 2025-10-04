package service

import (
	"context"

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

// CreateOrder создает новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, orderNumber string, userID uuid.UUID) (*domain.Order, error) {
	// Пытаемся создать заказ
	order, err := s.orderRepo.CreateOrder(ctx, orderNumber, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			// Заказ уже существует, проверим владельца
			existingOrders, getErr := s.orderRepo.GetOrdersByUserID(ctx, userID)
			if getErr != nil {
				return nil, getErr
			}

			// Ищем заказ среди заказов пользователя
			for _, o := range existingOrders {
				if o.Number == orderNumber {
					return &o, nil // Заказ принадлежит этому пользователю
				}
			}

			// Заказ принадлежит другому пользователю
			return nil, domain.ErrOrderBelongsToOtherUser
		}
		return nil, err
	}

	return order, nil
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
