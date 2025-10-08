package handlers

import (
	"errors"

	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/usecase"
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
	// if !util.ValidateLuhnNumber(numStr) {
	// 	c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Not valid luhn algorithm for number %s", numStr)})
	// }
	order, err := ah.uc.Order(c.Context(), numStr)

	if err != nil {
		zl.Log.Error("Internal error", zap.Error(err))
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(order)
}
func (ah *AccrualHandler) RegisterOrder(c fiber.Ctx) error {

	return c.SendStatus(fiber.StatusNotImplemented)

}
func (ah *AccrualHandler) RegisterGood(c fiber.Ctx) error {
	var rule domain.RewardRule
	if err := c.Bind().Body(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := ah.uc.RegisterGood(c.Context(), &rule)
	if errors.Is(err, repo.RuleAlreadyExistErr) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return err
}
func (ah *AccrualHandler) Ping(c fiber.Ctx) error {
	return ah.uc.Ping(c.RequestCtx())

}
