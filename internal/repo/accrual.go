package repo

import (
	"context"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccrualRepository struct {
	pool *pgxpool.Pool
}

func NewAccrualRepository(cfg *config.Accrual) *AccrualRepository {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}
	return &AccrualRepository{pool: pool}
}
