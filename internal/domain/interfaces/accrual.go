package interfaces

import (
	"context"
	"database/sql"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/handlers/requests"
)

type AccrualRepository interface {
	Order(ctx context.Context, number string) (*domain.AccrualOrder, error)
	RegisterOrder(ctx context.Context, order *requests.AccrualOrderRequest) error
	RegisterGood(ctx context.Context, rule *requests.RewardRuleRequest) error
	Ping(ctx context.Context) error
	StdDB() *sql.DB
	ProcessOrders(ctx context.Context, limit uint64) error
}
