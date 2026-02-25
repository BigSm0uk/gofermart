package service

import (
	"context"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/google/uuid"
)

// LoyaltyService представляет сервис для работы с баллами лояльности
type LoyaltyService struct {
	loyaltyRepo *repo.LoyaltyRepository
	orderRepo   *repo.OrderRepository
}

// NewLoyaltyService создает новый сервис баллов лояльности
func NewLoyaltyService(loyaltyRepo *repo.LoyaltyRepository, orderRepo *repo.OrderRepository) *LoyaltyService {
	return &LoyaltyService{
		loyaltyRepo: loyaltyRepo,
		orderRepo:   orderRepo,
	}
}

// GetUserBalance получает баланс пользователя
func (s *LoyaltyService) GetUserBalance(ctx context.Context, userID uuid.UUID) (*domain.BalanceResponse, error) {
	return s.loyaltyRepo.GetUserBalance(ctx, userID)
}

// WithdrawFunds списывает средства с баланса пользователя
func (s *LoyaltyService) WithdrawFunds(ctx context.Context, userID uuid.UUID, req *domain.WithdrawRequest) error {
	// Списываем средства
	return s.loyaltyRepo.WithdrawFunds(ctx, userID, req.Order, req.Sum)
}

// GetUserWithdrawals получает историю списаний пользователя
func (s *LoyaltyService) GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]domain.WithdrawalResponse, error) {
	withdrawals, err := s.loyaltyRepo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Преобразуем время в RFC3339 формат
	responses := make([]domain.WithdrawalResponse, len(withdrawals))
	for i, withdrawal := range withdrawals {
		responses[i] = domain.WithdrawalResponse{
			Order:       withdrawal.Order,
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt, // Уже в RFC3339 формате
		}
	}

	return responses, nil
}

// ProcessOrderAccrual обрабатывает начисление баллов за заказ (используется воркером)
func (s *LoyaltyService) ProcessOrderAccrual(ctx context.Context, userID uuid.UUID, orderNumber string, accrual float64) error {
	// Создаем операцию начисления
	_, err := s.loyaltyRepo.CreateOperation(ctx, userID, &orderNumber, domain.OperationTypeCredit, accrual)
	return err
}
