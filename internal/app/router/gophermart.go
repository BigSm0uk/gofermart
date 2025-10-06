package router

import (
	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/BigSm0uk/gofermart/internal/handlers"
	"github.com/BigSm0uk/gofermart/internal/handlers/middleware"
	"github.com/BigSm0uk/gofermart/internal/service"
	"github.com/gofiber/fiber/v3"
)

// SetupGophermartRoutes настраивает маршруты для gophermart
func SetupGophermartRoutes(app *fiber.App, userService *service.UserService, orderService *service.OrderService, loyaltyService *service.LoyaltyService, jwtManager *auth.JWTManager) {
	// Создаем хендлер
	handler := handlers.NewGophermartHandler(userService, orderService, loyaltyService)

	// Публичные маршруты (без аутентификации)
	api := app.Group("/api/user")
	api.Post("/register", handler.RegisterUser())
	api.Post("/login", handler.LoginUser())

	// Защищенные маршруты (требуют аутентификации)
	protected := api.Group("", middleware.AuthMiddleware(jwtManager))
	protected.Post("/orders", handler.CreateOrder())
	protected.Get("/orders", handler.GetUserOrders())
	protected.Get("/balance", handler.GetUserBalance())
	protected.Post("/balance/withdraw", handler.WithdrawFunds())
	protected.Get("/withdrawals", handler.GetUserWithdrawals())

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})
}