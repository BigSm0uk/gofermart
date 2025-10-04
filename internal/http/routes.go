package http

import (
	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes настраивает маршруты приложения
func SetupRoutes(app *fiber.App, handler *Handler, jwtManager *auth.JWTManager) {
	// Публичные маршруты (без аутентификации)
	api := app.Group("/api/user")
	api.Post("/register", handler.RegisterUser)
	api.Post("/login", handler.LoginUser)

	// Защищенные маршруты (требуют аутентификации)
	protected := api.Group("", AuthMiddleware(jwtManager))
	protected.Post("/orders", handler.CreateOrder)
	protected.Get("/orders", handler.GetUserOrders)
	protected.Get("/balance", handler.GetUserBalance)
	protected.Post("/balance/withdraw", handler.WithdrawFunds)
	protected.Get("/withdrawals", handler.GetUserWithdrawals)

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})
}
