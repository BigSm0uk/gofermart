package gophermart

import (
	"context"
	"fmt"
	"time"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/router"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/BigSm0uk/gofermart/internal/handlers/middleware"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Container struct {
	Server          *fiber.App
	Config          *config.Gophermart
	DB              *pgxpool.Pool
	UserRepo        *repo.UserRepository
	OrderRepo       *repo.OrderRepository
	LoyaltyRepo     *repo.LoyaltyRepository
	PasswordManager *auth.PasswordManager
	JWTManager      *auth.JWTManager
	UserService     *service.UserService
	OrderService    *service.OrderService
	LoyaltyService  *service.LoyaltyService
	Worker          *service.OrderProcessor
	Log             *zap.Logger
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) LoadConfig() (*Container, error) {
	cfg, err := config.InitGophermartConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	c.Config = cfg
	return c, nil
}

func (c *Container) LoadLogger() (*Container, error) {
	c.Log = zl.InitLogger(c.Config.Env)
	return c, nil
}

func (c *Container) LoadDatabase() (*Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, c.Config.Storage.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	c.DB = db
	c.Log.Info("Connected to database successfully")

	// Применяем миграции сразу после подключения
	if err := c.applyMigrations(); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	c.Log.Info("Database migrations applied successfully")

	return c, nil
}

// applyMigrations применяет миграции базы данных
func (c *Container) applyMigrations() error {
	// Получаем стандартное соединение для goose
	db := stdlib.OpenDBFromPool(c.DB)
	return Migrate(db)
}

func (c *Container) LoadRepositories() (*Container, error) {
	c.UserRepo = repo.NewUserRepository(c.DB)
	c.OrderRepo = repo.NewOrderRepository(c.DB)
	c.LoyaltyRepo = repo.NewLoyaltyRepository(c.DB)
	return c, nil
}

func (c *Container) LoadAuth() (*Container, error) {
	c.PasswordManager = auth.NewPasswordManager()
	c.JWTManager = auth.NewJWTManager(c.Config.JWTSecret, 24*time.Hour)
	return c, nil
}

func (c *Container) LoadServices() (*Container, error) {
	c.UserService = service.NewUserService(c.UserRepo, c.PasswordManager, c.JWTManager)
	c.OrderService = service.NewOrderService(c.OrderRepo)
	c.LoyaltyService = service.NewLoyaltyService(c.LoyaltyRepo, c.OrderRepo)
	return c, nil
}

func (c *Container) LoadWorker() (*Container, error) {
	c.Worker = service.NewOrderProcessor(c.OrderService, c.LoyaltyService, c.Config.WorkerInterval, c.Config.AccrualSystemAddress)
	return c, nil
}

func (c *Container) LoadServer() (*Container, error) {
	c.Server = fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Настраиваем middleware
	c.Server.Use(recover.New())
	c.Server.Use(cors.New())
	c.Server.Use(middleware.RequestIDMiddleware())

	// Настраиваем маршруты
	router.SetupGophermartRoutes(c.Server, c.UserService, c.OrderService, c.LoyaltyService, c.JWTManager, c.Log)

	return c, nil
}

func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

func Build() (*Container, error) {
	c := NewContainer()

	var err error
	if c, err = c.LoadConfig(); err != nil {
		return nil, err
	}
	if c, err = c.LoadLogger(); err != nil {
		return nil, err
	}
	if c, err = c.LoadDatabase(); err != nil {
		return nil, err
	}
	if c, err = c.LoadRepositories(); err != nil {
		return nil, err
	}
	if c, err = c.LoadAuth(); err != nil {
		return nil, err
	}
	if c, err = c.LoadServices(); err != nil {
		return nil, err
	}
	if c, err = c.LoadWorker(); err != nil {
		return nil, err
	}
	if c, err = c.LoadServer(); err != nil {
		return nil, err
	}

	return c, nil
}
