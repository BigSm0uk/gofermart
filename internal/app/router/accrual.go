package router

import (
	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/handlers"
	"github.com/BigSm0uk/gofermart/internal/handlers/middleware"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AccrualRouter struct {
	cfg    *config.Accrual
	server *fiber.App
	h      *handlers.AccrualHandler
	log    *zap.Logger
}

func NewAccrualRouter(cfg *config.Accrual, h *handlers.AccrualHandler, log *zap.Logger) *AccrualRouter {
	server := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	server.Use(func(c *fiber.Ctx) error {
		return middleware.LoggerMiddleware(c, log)
	})

	setupRoutes(server, h)

	return &AccrualRouter{cfg: cfg, server: server, h: h, log: log}
}
func (r *AccrualRouter) Listen() error {
	return r.server.Listen(r.cfg.RunAddress)
}
func (r *AccrualRouter) Shutdown() {
	r.server.Shutdown()
}

func setupRoutes(server *fiber.App, h *handlers.AccrualHandler) {
	server.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("pong")
	})
	server.Get("/api/orders/:number", h.Orders)
	server.Post("/api/orders", h.RegisterOrder)
	server.Post("/api/goods", h.RegisterGood)
}
