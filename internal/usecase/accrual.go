package usecase

import (
	"context"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/domain/interfaces"
)

type AccrualUsecase struct {
	repo interfaces.AccrualRepository
}

func NewAccrualUsecase(r interfaces.AccrualRepository) *AccrualUsecase {
	return &AccrualUsecase{repo: r}
}
func (a *AccrualUsecase) Order(ctx context.Context, number string) (*domain.AccrualOrder, error) {
	return a.repo.Order(ctx, number)
}
func (a *AccrualUsecase) RegisterOrder(ctx context.Context, order *domain.AccrualOrder) error {
	return a.repo.RegisterOrder(ctx, order)
}
func (a *AccrualUsecase) RegisterGood(ctx context.Context, rule *domain.RewardRule) error {
	err := a.repo.RegisterGood(ctx, rule)

	return err
}
func (a *AccrualUsecase) Ping(ctx context.Context) error {
	return a.repo.Ping(ctx)
}
