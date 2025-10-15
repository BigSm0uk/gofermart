package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LoyaltyRepository представляет репозиторий для работы с баллами лояльности
type LoyaltyRepository struct {
	db *pgxpool.Pool
}

// NewLoyaltyRepository создает новый репозиторий баллов лояльности
func NewLoyaltyRepository(db *pgxpool.Pool) *LoyaltyRepository {
	return &LoyaltyRepository{db: db}
}

// CreateOperation создает операцию с баллами лояльности
func (r *LoyaltyRepository) CreateOperation(ctx context.Context, userID uuid.UUID, orderNumber *string, operationType domain.OperationType, amount float64) (*domain.LoyaltyOperation, error) {
	query := `
		INSERT INTO loyalty_operations (user_id, order_number, operation_type, amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, order_number, operation_type, amount, processed_at, created_at`

	var operation domain.LoyaltyOperation
	err := r.db.QueryRow(ctx, query, userID, orderNumber, operationType, amount).Scan(
		&operation.ID,
		&operation.UserID,
		&operation.OrderNumber,
		&operation.OperationType,
		&operation.Amount,
		&operation.ProcessedAt,
		&operation.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create loyalty operation: %w", err)
	}

	return &operation, nil
}

// GetUserBalance вычисляет текущий баланс пользователя
func (r *LoyaltyRepository) GetUserBalance(ctx context.Context, userID uuid.UUID) (*domain.BalanceResponse, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN operation_type = 'CREDIT' THEN amount ELSE -amount END), 0) as current,
			COALESCE(SUM(CASE WHEN operation_type = 'DEBIT' THEN amount ELSE 0 END), 0) as withdrawn
		FROM loyalty_operations
		WHERE user_id = $1`

	var balance domain.BalanceResponse
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&balance.Current,
		&balance.Withdrawn,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	return &balance, nil
}

// GetUserWithdrawals получает все списания пользователя, отсортированные по времени
func (r *LoyaltyRepository) GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]domain.WithdrawalResponse, error) {
	query := `
		SELECT order_number, amount, processed_at
		FROM loyalty_operations
		WHERE user_id = $1 AND operation_type = 'DEBIT'
		ORDER BY processed_at ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user withdrawals: %w", err)
	}
	defer rows.Close()

	var withdrawals []domain.WithdrawalResponse
	for rows.Next() {
		var withdrawal domain.WithdrawalResponse
		var orderNumber *string
		var processedAt interface{}

		err := rows.Scan(&orderNumber, &withdrawal.Sum, &processedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan withdrawal: %w", err)
		}

		if orderNumber != nil {
			withdrawal.Order = *orderNumber
		}

		// Преобразуем время в строку RFC3339
		if processedAt != nil {
			if timeValue, ok := processedAt.(time.Time); ok {
				withdrawal.ProcessedAt = timeValue.Format("2006-01-02T15:04:05Z07:00")
			}
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating withdrawals: %w", err)
	}

	return withdrawals, nil
}

// WithdrawFunds списывает средства с баланса пользователя в транзакции
//
// Проблема конкурентности:
//   - READ COMMITTED (по умолчанию): между чтением баланса и записью операции
//     другой запрос может изменить баланс → возможен овердрафт
//   - REPEATABLE READ: гарантирует консистентность данных в рамках транзакции,
//     при конфликте сериализации - ошибка, которую обрабатываем повторной попыткой
//   - SERIALIZABLE: максимальная изоляция, но избыточно для данной задачи
//
// Решение: REPEATABLE READ + retry при ошибках сериализации
func (r *LoyaltyRepository) WithdrawFunds(ctx context.Context, userID uuid.UUID, orderNumber string, amount float64) error {
	// Максимальное количество попыток при ошибке сериализации
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := r.withdrawFundsWithRetry(ctx, userID, orderNumber, amount)
		if err == nil {
			return nil // Успешно
		}

		// Проверяем, является ли это ошибкой сериализации
		if !isSerializationError(err) {
			return err // Не ошибка сериализации, возвращаем как есть
		}

		// Ошибка сериализации - делаем паузу и повторяем
		if attempt < maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * 10 * time.Millisecond)
		}
	}

	return fmt.Errorf("withdrawal failed after %d attempts due to concurrent modifications", maxRetries)
}

// withdrawFundsWithRetry выполняет списание средств в транзакции с REPEATABLE READ
func (r *LoyaltyRepository) withdrawFundsWithRetry(ctx context.Context, userID uuid.UUID, orderNumber string, amount float64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Устанавливаем уровень изоляции REPEATABLE READ
	_, err = tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		return fmt.Errorf("failed to set transaction isolation level: %w", err)
	}

	// Проверяем текущий баланс
	balanceQuery := `
		SELECT COALESCE(SUM(CASE WHEN operation_type = 'CREDIT' THEN amount ELSE -amount END), 0)
		FROM loyalty_operations
		WHERE user_id = $1`

	var currentBalance float64
	err = tx.QueryRow(ctx, balanceQuery, userID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	if currentBalance < amount {
		return domain.ErrInsufficientFunds
	}

	// Создаем операцию списания
	insertQuery := `
		INSERT INTO loyalty_operations (user_id, order_number, operation_type, amount)
		VALUES ($1, $2, 'DEBIT', $3)`

	_, err = tx.Exec(ctx, insertQuery, userID, orderNumber, amount)
	if err != nil {
		return fmt.Errorf("failed to create withdrawal operation: %w", err)
	}

	return tx.Commit(ctx)
}

// isSerializationError проверяет, является ли ошибка ошибкой сериализации
func isSerializationError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "serialization failure") ||
		strings.Contains(errStr, "could not serialize access")
}
