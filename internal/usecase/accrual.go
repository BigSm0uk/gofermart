package usecase

import (
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/repo"
)

type AccrualUsecase struct {
	repo *repo.AccrualRepository
}

func NewAccrualUsecase(r *repo.AccrualRepository) *AccrualUsecase {
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
