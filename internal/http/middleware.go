package http

import (
	"log"
	"strings"
	"time"

	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// AuthMiddleware создает middleware для аутентификации
func AuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c fiber.Ctx) error {
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

// RequestIDMiddleware создает middleware для добавления request ID
func RequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// LoggingMiddleware создает middleware для логирования запросов
func LoggingMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		requestID := c.Locals("requestID").(string)

		// Логируем начало запроса
		log.Printf("[%s] %s %s started", requestID, c.Method(), c.Path())

		// Выполняем следующий хендлер
		err := c.Next()

		// Логируем результат
		duration := time.Since(start)
		status := c.Response().StatusCode()

		if err != nil {
			log.Printf("[%s] %s %s failed: %v (duration: %v)", requestID, c.Method(), c.Path(), err, duration)
		} else {
			log.Printf("[%s] %s %s completed with status %d (duration: %v)", requestID, c.Method(), c.Path(), status, duration)
		}

		return err
	}
}

// ErrorHandler обрабатывает ошибки
func ErrorHandler(c fiber.Ctx, err error) error {
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

// generateRequestID генерирует уникальный ID запроса
func generateRequestID() string {
	return uuid.New().String()[:8]
}
