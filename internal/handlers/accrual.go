package handlers

import (
	"errors"

	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/handlers/requests"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/usecase"
	"github.com/BigSm0uk/gofermart/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type AccrualHandler struct {
	uc *usecase.AccrualUsecase
}

func NewAccrualHandler(uc *usecase.AccrualUsecase) *AccrualHandler {
	return &AccrualHandler{uc: uc}
}

func (ah *AccrualHandler) Orders(c fiber.Ctx) error {
	numStr := c.Params("number")
	isValid := utils.ValidateLuhn(numStr)
	if !isValid {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Invalid order number format"})
	}
	order, err := ah.uc.Order(c.Context(), numStr)

	if err != nil {
		zl.Log.Error("Internal error", zap.Error(err))
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(order)
}
func (ah *AccrualHandler) RegisterOrder(c fiber.Ctx) error {
	var order requests.AccrualOrderRequest
	if err := c.Bind().Body(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	isValid := utils.ValidateLuhn(order.Order)
	if !isValid {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Invalid order number format"})
	}
	err := ah.uc.RegisterOrder(c.Context(), &order)
	if err != nil {
		if errors.Is(err, repo.ErrOrderAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusAccepted)

}
func (ah *AccrualHandler) RegisterGood(c fiber.Ctx) error {
	var rule requests.RewardRuleRequest
	if err := c.Bind().Body(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := ah.uc.RegisterGood(c.Context(), &rule)
	if errors.Is(err, repo.ErrRuleAlreadyExist) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return err
}
func (ah *AccrualHandler) Ping(c fiber.Ctx) error {
	return ah.uc.Ping(c.RequestCtx())

}
