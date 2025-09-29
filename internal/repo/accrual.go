package repo

import (
	"context"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/domain/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AccrualRepository struct {
	pool *pgxpool.Pool
}

var _ interfaces.AccrualRepository = &AccrualRepository{}

func NewAccrualRepository(cfg *config.Accrual) *AccrualRepository {
	zl.Log.Info("NewAccrualRepository", zap.String("databaseURI", cfg.Storage.DatabaseURI))
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
