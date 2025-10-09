package accrual

import (
	"context"
	"time"

	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/domain/interfaces"
	"go.uber.org/zap"
)

type AccrualProcessor struct {
	repo  interfaces.AccrualRepository
	ttl   time.Duration
	limit uint64
}

func NewAccrualProcessor(repo interfaces.AccrualRepository, ttl time.Duration, limit uint64) *AccrualProcessor {
	return &AccrualProcessor{repo: repo, ttl: ttl, limit: limit}
}
func (p *AccrualProcessor) Run(ctx context.Context) {
	ticker := time.NewTicker(p.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processOrders(ctx)
		}
	}
}
func (p *AccrualProcessor) processOrders(ctx context.Context) {
	err := p.repo.ProcessOrders(ctx, p.limit)
	if err != nil {
		zl.Log.Error("failed to get orders for processing", zap.Error(err))
		return
	}
}
