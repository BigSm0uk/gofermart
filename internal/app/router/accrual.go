package router

import (
	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type AccrualRouter struct {
	cfg    *config.Accrual
	server *fiber.App
}

func NewAccrualRouter(cfg *config.Accrual) *AccrualRouter {
	server := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	return &AccrualRouter{cfg: cfg, server: server}
}
