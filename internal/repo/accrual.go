package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/domain/interfaces"
	"github.com/BigSm0uk/gofermart/internal/handlers/requests"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var (
	ErrRuleAlreadyExist   = errors.New("rule already exist")
	ErrOrderAlreadyExists = errors.New("order already exists")
)

type AccrualRepository struct {
	pool *pgxpool.Pool
}

var _ interfaces.AccrualRepository = &AccrualRepository{}

func NewAccrualRepository(cfg *config.Accrual) *AccrualRepository {
	pCfg, err := pgxpool.ParseConfig(cfg.Storage.DatabaseURI)
	pCfg.MinConns = cfg.Storage.MinPoolSize
	pCfg.MaxConns = cfg.Storage.MaxPoolSize
	pCfg.MaxConnLifetime = cfg.Storage.ConnectionLifetime
	if err != nil {
		panic(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), pCfg)
	if err != nil {
		panic(err)
	}
	return &AccrualRepository{pool: pool}
}
func (a *AccrualRepository) Order(ctx context.Context, number string) (*domain.AccrualOrder, error) {
	sql, args, err := sq.
		Select("id, order_number, status, accrual, created_at").
		From("accrual_orders").
		Where(sq.Eq{"order_number": number}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var order domain.AccrualOrder
	err = a.pool.QueryRow(ctx, sql, args...).Scan(&order.ID, &order.Order, &order.Status, &order.Accrual, &order.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &order, nil
}
func (a *AccrualRepository) RegisterOrder(ctx context.Context, order *requests.AccrualOrderRequest) error {
	tx, err := a.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql, args, err := sq.
		Insert("accrual_orders").
		Columns("order_number").
		Values(order.Order).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	var id int64
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		zl.Log.Error("register order error", zap.Error(err))
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrOrderAlreadyExists
		}
		return err
	}

	for _, good := range order.Goods {
		sql, args, err = sq.
			Insert("accrual_order_goods").
			Columns("order_id", "description", "price").
			Values(id, good.Description, good.Price).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			zl.Log.Error("register good error", zap.Error(err))
			return err
		}
		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			zl.Log.Error("register good error", zap.Error(err))
			return err
		}
	}
	return tx.Commit(ctx)
}

func (a *AccrualRepository) RegisterGood(ctx context.Context, rule *requests.RewardRuleRequest) error {
	sql, args, err := sq.
		Insert("reward_rules").
		Columns("match", "reward", "reward_type").
		Values(rule.Match, rule.Reward, rule.RewardType).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = a.pool.Exec(ctx, sql, args...)
	if err != nil {
		zl.Log.Error("register good error", zap.Error(err))
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrRuleAlreadyExist
		}
		return err
	}

	return nil

}

func (a *AccrualRepository) Ping(ctx context.Context) error {
	return a.pool.Ping(ctx)
}
func (a *AccrualRepository) StdDB() *sql.DB {
	return stdlib.OpenDBFromPool(a.pool)
}
func (a *AccrualRepository) ProcessOrders(ctx context.Context, limit uint64) error {
	// Шаг 1: короткая транзакция — выбрать партию с блокировкой и пометить PROCESSING
	tx, err := a.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql, args, err := sq.
		Select("id, order_number").
		From("accrual_orders").
		Where("status in ($1, $2)", domain.AccrualStatusRegistered, domain.AccrualStatusProcessing).
		OrderBy("created_at ASC").
		Suffix("FOR UPDATE SKIP LOCKED").
		Limit(limit).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var orders []domain.AccrualOrder
	for rows.Next() {
		var order domain.AccrualOrder
		if err := rows.Scan(&order.ID, &order.Order); err != nil {
			return fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating orders: %w", err)
	}

	for _, order := range orders {
		sql, args, err = sq.
			Update("accrual_orders").
			Set("status", domain.AccrualStatusProcessing).
			Where(sq.Eq{"id": order.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
		if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Шаг 2: обработать каждый заказ в отдельной транзакции
	for _, order := range orders {
		if err := a.processSingleOrder(ctx, order.ID, order.Order); err != nil {
			zl.Log.Error("failed to process order", zap.Error(err), zap.String("order_number", order.Order))
		}
	}
	return nil
}

func (a *AccrualRepository) processSingleOrder(ctx context.Context, orderID int64, orderNumber string) error {
	tx, err := a.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	accrual, status := a.calculateAccrual(ctx, tx, orderID)

	sql, args, err := sq.
		Update("accrual_orders").
		Set("status", status).
		Set("accrual", accrual).
		Where(sq.Eq{"id": orderID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update accrual sql: %w", err)
	}
	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed to update order %s: %w", orderNumber, err)
	}
	zl.Log.Info("order processed", zap.String("order_number", orderNumber), zap.Float64("accrual", accrual), zap.String("status", string(status)))
	return tx.Commit(ctx)
}
func (a *AccrualRepository) calculateAccrual(ctx context.Context, tx pgx.Tx, orderID int64) (float64, domain.AccrualOrderStatus) {
	sql, args, err := sq.Select("COALESCE(SUM(CASE rr.reward_type WHEN '%' THEN aog.price * rr.reward / 100.0 WHEN 'pt' THEN rr.reward END), 0) AS accrual").
		From("accrual_order_goods aog").
		Join("reward_rules rr ON aog.description ILIKE '%' || rr.match || '%' ").
		Where(sq.Eq{"aog.order_id": orderID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, domain.AccrualStatusInvalid
	}
	var accrual float64
	err = tx.QueryRow(ctx, sql, args...).Scan(&accrual)
	if err != nil {
		return 0, domain.AccrualStatusInvalid
	}
	if accrual == 0 {
		return 0, domain.AccrualStatusInvalid
	}
	return accrual, domain.AccrualStatusProcessed
}
