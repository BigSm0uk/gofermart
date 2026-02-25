package handlers

import (
	"errors"

	"github.com/BigSm0uk/gofermart/internal/handlers/requests"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/usecase"
	"github.com/BigSm0uk/gofermart/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AccrualHandler struct {
	uc        *usecase.AccrualUsecase
	validator StructValidator
	log       *zap.Logger
}

func NewAccrualHandler(uc *usecase.AccrualUsecase, log *zap.Logger) *AccrualHandler {
	return &AccrualHandler{
		uc:        uc,
		validator: NewStructValidator(),
		log:       log,
	}
}

func (ah *AccrualHandler) Orders(c *fiber.Ctx) error {
	numStr := c.Params("number")
	isValid := utils.ValidateLuhn(numStr)
	if !isValid {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Invalid order number format"})
	}
	order, err := ah.uc.Order(c.Context(), numStr)

	if err != nil {
		ah.log.Error("Internal error", zap.Error(err))
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(order)
}
func (ah *AccrualHandler) RegisterOrder(c *fiber.Ctx) error {
	var order requests.AccrualOrderRequest
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Валидация структуры
	if err := ah.validator.Validate(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed: " + err.Error()})
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
	return c.Status(fiber.StatusAccepted).Send(nil)

}
func (ah *AccrualHandler) RegisterGood(c *fiber.Ctx) error {
	var rule requests.RewardRuleRequest
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Валидация структуры
	if err := ah.validator.Validate(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed: " + err.Error()})
	}

	err := ah.uc.RegisterGood(c.Context(), &rule)
	if errors.Is(err, repo.ErrRuleAlreadyExist) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return err
}
func (ah *AccrualHandler) Ping(c *fiber.Ctx) error {
	return ah.uc.Ping(c.Context())

}
