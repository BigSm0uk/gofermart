package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware создает middleware для добавления request ID
func RequestIDMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// generateRequestID генерирует уникальный ID запроса
func generateRequestID() string {
	return uuid.New().String()[:8]
}
