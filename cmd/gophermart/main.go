package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/BigSm0uk/gofermart/internal/http"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/service"
	"github.com/BigSm0uk/gofermart/internal/worker"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	config := loadConfig()

	// Подключаемся к базе данных
	db, err := connectDB(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Создаем репозитории
	userRepo := repo.NewUserRepository(db)
	orderRepo := repo.NewOrderRepository(db)
	loyaltyRepo := repo.NewLoyaltyRepository(db)

	// Создаем менеджеры аутентификации
	passwordManager := auth.NewPasswordManager()
	jwtManager := auth.NewJWTManager(config.JWTSecret, 24*time.Hour) // Токен действует 24 часа

	// Создаем сервисы
	userService := service.NewUserService(userRepo, passwordManager, jwtManager)
	orderService := service.NewOrderService(orderRepo)
	loyaltyService := service.NewLoyaltyService(loyaltyRepo, orderRepo)

	// Создаем HTTP хендлеры
	handler := http.NewHandler(userService, orderService, loyaltyService)

	// Создаем Fiber приложение
	app := fiber.New(fiber.Config{
		ErrorHandler: http.ErrorHandler,
	})

	// Настраиваем middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(http.RequestIDMiddleware())
	app.Use(http.LoggingMiddleware())

	// Настраиваем маршруты
	http.SetupRoutes(app, handler, jwtManager)

	// Создаем и запускаем воркер
	processor := worker.NewOrderProcessor(orderService, loyaltyService, config.WorkerInterval, config.AccrualSystemAddress)

	// Запускаем воркер в отдельной горутине
	ctx, cancel := context.WithCancel(context.Background())
	go processor.Run(ctx)

	// Запускаем HTTP сервер
	go func() {
		log.Printf("Starting server on %s", config.HTTPAddr)
		if err := app.Listen(config.HTTPAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Останавливаем воркер
	cancel()

	// Останавливаем сервер
	if err := app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// Config представляет конфигурацию приложения
type Config struct {
	HTTPAddr             string
	DatabaseURL          string
	JWTSecret            string
	WorkerInterval       time.Duration
	AccrualSystemAddress string
}

// loadConfig загружает конфигурацию из переменных окружения
func loadConfig() *Config {
	httpAddr := getEnv("HTTP_ADDR", ":8080")
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	workerIntervalStr := getEnv("WORKER_INTERVAL", "5s")
	accrualSystemAddress := getEnv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8081")

	workerInterval, err := time.ParseDuration(workerIntervalStr)
	if err != nil {
		log.Fatalf("Invalid WORKER_INTERVAL format: %v", err)
	}

	return &Config{
		HTTPAddr:             httpAddr,
		DatabaseURL:          databaseURL,
		JWTSecret:            jwtSecret,
		WorkerInterval:       workerInterval,
		AccrualSystemAddress: accrualSystemAddress,
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// connectDB подключается к базе данных PostgreSQL
func connectDB(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database successfully")
	return db, nil
}
