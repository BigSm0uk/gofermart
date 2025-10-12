package middleware

import (
	"strings"

	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware создает middleware для аутентификации
func AuthMiddleware(jwtManager *auth.JWTManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Получаем токен из заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Authorization header required"})
		}

		// Проверяем формат Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Парсим userID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID in token"})
		}

		// Сохраняем userID в контексте
		c.Locals("userID", userID)
		c.Locals("userLogin", claims.Login)

		return c.Next()
	}
}
