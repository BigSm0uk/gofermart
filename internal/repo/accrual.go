package repo

import (
	"context"
	"database/sql"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/domain/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type AccrualRepository struct {
	pool *pgxpool.Pool
}

var _ interfaces.AccrualRepository = &AccrualRepository{}

func NewAccrualRepository(cfg *config.Accrual) *AccrualRepository {
	pCfg, err := pgxpool.ParseConfig(cfg.Storage.DatabaseURI)
	if err != nil {
		panic(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), pCfg)
	if err != nil {
		panic(err)
	}
	return &AccrualRepository{pool: pool}
}
func (a *AccrualRepository) Orders(number int) (*domain.AccrualOrder, error) {
	panic("unimplemented")
}
func (a *AccrualRepository) RegisterGood(good *domain.AccrualOrderGood) error {
	panic("unimplemented")
}

func (a *AccrualRepository) RegisterOrder(order *domain.AccrualOrder) error {
	panic("unimplemented")
}

func (a *AccrualRepository) Ping(ctx context.Context) error {
	return a.pool.Ping(ctx)
}
func (a *AccrualRepository) GetStdDB() *sql.DB {
	return stdlib.OpenDBFromPool(a.pool)
}
