package repo

import (
	"context"
	"fmt"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepository представляет репозиторий для работы с заказами
type OrderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository создает новый репозиторий заказов
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder создает новый заказ или возвращает существующий
// Возвращает order, isNew (true если заказ был создан, false если уже существовал), error
func (r *OrderRepository) CreateOrder(ctx context.Context, number string, userID uuid.UUID) (*domain.Order, bool, error) {
	query := `
		INSERT INTO orders (number, user_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, number, user_id, status, accrual, uploaded_at, updated_at`

	var order domain.Order
	err := r.db.QueryRow(ctx, query, number, userID, domain.OrderStatusNew).Scan(
		&order.ID,
		&order.Number,
		&order.UserID,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if isUniqueViolation(err) {
			// Заказ уже существует, получаем существующий заказ
			existingOrder, getErr := r.getOrderByNumber(ctx, number)
			if getErr != nil {
				return nil, false, fmt.Errorf("failed to get existing order: %w", getErr)
			}
			return existingOrder, false, nil // Заказ уже существовал
		}
		return nil, false, fmt.Errorf("failed to create order: %w", err)
	}

	return &order, true, nil // Новый заказ
}

// getOrderByNumber получает заказ по номеру
func (r *OrderRepository) getOrderByNumber(ctx context.Context, number string) (*domain.Order, error) {
	query := `
		SELECT id, number, user_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE number = $1`

	var order domain.Order
	err := r.db.QueryRow(ctx, query, number).Scan(
		&order.ID,
		&order.Number,
		&order.UserID,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get order by number: %w", err)
	}

	return &order, nil
}

// GetOrdersByUserID получает все заказы пользователя, отсортированные по времени загрузки
func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	query := `
		SELECT id, number, user_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user ID: %w", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.Number,
			&order.UserID,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

// UpdateOrderStatus обновляет статус заказа
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, number string, status domain.OrderStatus, accrual *float64) error {
	var query string
	var args []interface{}

	if accrual != nil {
		query = `
			UPDATE orders 
			SET status = $1, accrual = $2, updated_at = NOW()
			WHERE number = $3`
		args = []interface{}{status, *accrual, number}
	} else {
		query = `
			UPDATE orders 
			SET status = $1, updated_at = NOW()
			WHERE number = $2`
		args = []interface{}{status, number}
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetOrdersForProcessing получает заказы для обработки воркером
func (r *OrderRepository) GetOrdersForProcessing(ctx context.Context, limit int) ([]domain.Order, error) {
	query := `
		SELECT id, number, user_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE status IN ($1, $2)
		ORDER BY uploaded_at ASC
		LIMIT $3`

	rows, err := r.db.Query(ctx, query, domain.OrderStatusNew, domain.OrderStatusProcessing, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for processing: %w", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.Number,
			&order.UserID,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}
