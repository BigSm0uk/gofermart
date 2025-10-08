package repo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/domain/interfaces"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var (
	RuleAlreadyExistErr = errors.New("rule already exist")
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
func (a *AccrualRepository) RegisterOrder(ctx context.Context, order *domain.AccrualOrder) error {
	panic("unimplemented")
}

func (a *AccrualRepository) RegisterGood(ctx context.Context, rule *domain.RewardRule) error {
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
			return RuleAlreadyExistErr
		}
		return err
	}

	return nil

}

func (a *AccrualRepository) Ping(ctx context.Context) error {
	return a.pool.Ping(ctx)
}
func (a *AccrualRepository) GetStdDB() *sql.DB {
	return stdlib.OpenDBFromPool(a.pool)
}
