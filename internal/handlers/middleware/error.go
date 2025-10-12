package middleware

import (
	"log"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler обрабатывает ошибки
func ErrorHandler(c *fiber.Ctx, err error) error {
	requestID := c.Locals("requestID")
	if requestID == nil {
		requestID = "unknown"
	}

	log.Printf("[%s] Error: %v", requestID, err)

	// Определяем статус код на основе типа ошибки
	status := 500
	switch err {
	case domain.ErrNotFound:
		status = 404
	case domain.ErrInvalidCredentials:
		status = 401
	case domain.ErrUserAlreadyExists:
		status = 409
	case domain.ErrOrderAlreadyExists:
		status = 200 // Заказ уже существует у этого пользователя
	case domain.ErrOrderBelongsToOtherUser:
		status = 409
	case domain.ErrInvalidOrderNumber:
		status = 422
	case domain.ErrInsufficientFunds:
		status = 402
	case domain.ErrUnauthorized:
		status = 401
	}

	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}
