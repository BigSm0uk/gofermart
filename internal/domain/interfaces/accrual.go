package interfaces

import "github.com/BigSm0uk/gofermart/internal/domain"

type AccrualRepository interface {
	Orders(number int) (*domain.AccrualOrder, error)
	RegisterOrder(order *domain.AccrualOrder) error
	RegisterGood(good *domain.AccrualOrderGood) error
}
