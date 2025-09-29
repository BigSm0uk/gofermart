package handlers

import (
	"github.com/BigSm0uk/gofermart/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type AccrualHandler struct {
	uc *usecase.AccrualUsecase
}

func NewAccrualHandler(uc *usecase.AccrualUsecase) *AccrualHandler {
	return &AccrualHandler{uc: uc}
}

func (ah *AccrualHandler) Orders() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}
func (ah *AccrualHandler) RegisterOrder() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}
func (ah *AccrualHandler) RegisterGood() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}
func (ah *AccrualHandler) Ping() fiber.Handler {
	return func(c fiber.Ctx) error {
		return ah.uc.Ping(c.RequestCtx())
	}
}
