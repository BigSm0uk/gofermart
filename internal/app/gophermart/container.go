package gophermart

import (
	"context"
	"fmt"
	"time"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/router"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/BigSm0uk/gofermart/internal/handlers"
	"github.com/BigSm0uk/gofermart/internal/handlers/middleware"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
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
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) LoadConfig() *Container {
	config, err := config.InitGophermartConfig()
	if err != nil {
		panic(err)
	}
	c.Config = config
	return c
}

func (c *Container) LoadLogger() *Container {
	zl.InitLogger(c.Config.Env)
	return c
}

func (c *Container) LoadDatabase() *Container {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, c.Config.DatabaseURI)
	if err != nil {
		panic(fmt.Errorf("failed to create connection pool: %w", err))
	}

	// Проверяем соединение
	if err := db.Ping(ctx); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	c.DB = db
	zl.Log.Info("Connected to database successfully")
	return c
}

func (c *Container) LoadRepositories() *Container {
	c.UserRepo = repo.NewUserRepository(c.DB)
	c.OrderRepo = repo.NewOrderRepository(c.DB)
	c.LoyaltyRepo = repo.NewLoyaltyRepository(c.DB)
	return c
}

func (c *Container) LoadAuth() *Container {
	c.PasswordManager = auth.NewPasswordManager()
	c.JWTManager = auth.NewJWTManager(c.Config.JWTSecret, 24*time.Hour)
	return c
}

func (c *Container) LoadServices() *Container {
	c.UserService = service.NewUserService(c.UserRepo, c.PasswordManager, c.JWTManager)
	c.OrderService = service.NewOrderService(c.OrderRepo)
	c.LoyaltyService = service.NewLoyaltyService(c.LoyaltyRepo, c.OrderRepo)
	return c
}

func (c *Container) LoadWorker() *Container {
	c.Worker = service.NewOrderProcessor(c.OrderService, c.LoyaltyService, c.Config.WorkerInterval, c.Config.AccrualSystemAddress)
	return c
}

func (c *Container) LoadServer() *Container {
	// Создаем валидатор
	validator := handlers.NewStructValidator()

	c.Server = fiber.New(fiber.Config{
		ErrorHandler:     middleware.ErrorHandler,
		StructValidator:  validator,
	})

	// Настраиваем middleware
	c.Server.Use(recover.New())
	c.Server.Use(cors.New())
	c.Server.Use(middleware.RequestIDMiddleware())

	// Настраиваем маршруты
	router.SetupGophermartRoutes(c.Server, c.UserService, c.OrderService, c.LoyaltyService, c.JWTManager)

	return c
}

func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

func MustBuild() *Container {
	return NewContainer().
		LoadConfig().
		LoadLogger().
		LoadDatabase().
		LoadRepositories().
		LoadAuth().
		LoadServices().
		LoadWorker().
		LoadServer()
}