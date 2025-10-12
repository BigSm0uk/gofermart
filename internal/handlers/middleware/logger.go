package middleware

import (
	"time"

	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	stop := time.Now()

	latency := stop.Sub(start)
	status := c.Response().StatusCode()
	method := c.Method()
	path := c.Path()
	ip := c.IP()

	fields := []zap.Field{
		zap.Int("status", status),
		zap.String("method", method),
		zap.String("path", path),
		zap.String("ip", ip),
		zap.Duration("latency", latency),
	}

	if err != nil {
		zl.Log.Error("handle request with error", append(fields, zap.Error(err))...)
	} else {
		zl.Log.Info("handle request", fields...)
	}

	return err
}
