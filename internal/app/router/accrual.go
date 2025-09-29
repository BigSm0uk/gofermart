package router

import (
	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/handlers"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type AccrualRouter struct {
	cfg    *config.Accrual
	server *fiber.App
	h      *handlers.AccrualHandler
}

func NewAccrualRouter(cfg *config.Accrual, h *handlers.AccrualHandler) *AccrualRouter {
	server := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	setupRoutes(server, h)

	return &AccrualRouter{cfg: cfg, server: server, h: h}
}
func (r *AccrualRouter) Listen() error {
	return r.server.Listen(r.cfg.RunAddress)
}
func (r *AccrualRouter) Shutdown() {
	r.server.Shutdown()
}

func setupRoutes(server *fiber.App, h *handlers.AccrualHandler) {
	server.Get("/api/ping", func(c fiber.Ctx) error {
		c.SendString("pong")
		return c.SendStatus(fiber.StatusOK)
	})
	server.Get("/api/orders/:number", h.Orders())
	server.Post("/api/orders", h.RegisterOrder())
	server.Post("/api/goods", h.RegisterGood())
}
