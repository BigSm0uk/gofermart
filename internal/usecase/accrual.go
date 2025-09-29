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
func (a *AccrualUsecase) Orders(number int) (*domain.AccrualOrder, error) {
	return a.repo.Orders(number)
}
func (a *AccrualUsecase) RegisterOrder(order *domain.AccrualOrder) error {
	return a.repo.RegisterOrder(order)
}
func (a *AccrualUsecase) RegisterGood(good *domain.AccrualOrderGood) error {
	return a.repo.RegisterGood(good)
}
func (a *AccrualUsecase) Ping(ctx context.Context) error {
	return a.repo.Ping(ctx)
}
