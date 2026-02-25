package handlers

import (
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/service"
	"github.com/BigSm0uk/gofermart/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GophermartHandler представляет HTTP хендлеры для gophermart
type GophermartHandler struct {
	userService    *service.UserService
	orderService   *service.OrderService
	loyaltyService *service.LoyaltyService
	validator      StructValidator
	log            *zap.Logger
}

// NewGophermartHandler создает новый хендлер
func NewGophermartHandler(userService *service.UserService, orderService *service.OrderService, loyaltyService *service.LoyaltyService, log *zap.Logger) *GophermartHandler {
	return &GophermartHandler{
		userService:    userService,
		orderService:   orderService,
		loyaltyService: loyaltyService,
		validator:      NewStructValidator(),
		log:            log,
	}
}

// RegisterUser регистрирует нового пользователя
func (h *GophermartHandler) RegisterUser() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var req domain.UserCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Валидация структуры
		if err := h.validator.Validate(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Validation failed: " + err.Error()})
		}

		user, token, err := h.userService.RegisterUser(c.Context(), &req)
		if err != nil {
			return err
		}

		h.log.Info("User registered successfully", zap.String("login", user.Login))

		// Устанавливаем токен в заголовок Authorization
		c.Set("Authorization", "Bearer "+token)

		return c.Status(200).JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":    user.ID,
				"login": user.Login,
			},
		})
	}
}

// LoginUser аутентифицирует пользователя
func (h *GophermartHandler) LoginUser() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var req domain.UserLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Валидация структуры
		if err := h.validator.Validate(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Validation failed: " + err.Error()})
		}

		user, token, err := h.userService.LoginUser(c.Context(), &req)
		if err != nil {
			return err
		}

		h.log.Info("User logged in successfully", zap.String("login", user.Login))

		// Устанавливаем токен в заголовок Authorization
		c.Set("Authorization", "Bearer "+token)

		return c.Status(200).JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":    user.ID,
				"login": user.Login,
			},
		})
	}
}

// CreateOrder создает новый заказ
func (h *GophermartHandler) CreateOrder() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)
		orderNumber := string(c.Body())

		if len(orderNumber) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Order number is required"})
		}

		// Проверка номера заказа по алгоритму Луна
		if !utils.ValidateLuhn(orderNumber) {
			return c.Status(422).JSON(fiber.Map{"error": "Invalid order number format"})
		}

		order, isNew, err := h.orderService.CreateOrder(c.Context(), orderNumber, userID)
		if err != nil {
			return err
		}

		// Определяем статус ответа
		status := 202 // Новый заказ
		if !isNew {
			status = 200 // Заказ уже существовал
		}

		h.log.Info("Order created", zap.String("order", order.Number), zap.String("userID", userID.String()), zap.Bool("isNew", isNew))

		return c.Status(status).JSON(fiber.Map{
			"message": "Order processed",
			"order": fiber.Map{
				"number": order.Number,
				"status": order.Status,
			},
		})
	}
}

// GetUserOrders получает заказы пользователя
func (h *GophermartHandler) GetUserOrders() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		orders, err := h.orderService.GetUserOrders(c.Context(), userID)
		if err != nil {
			return err
		}

		if len(orders) == 0 {
			return c.Status(204).Send(nil)
		}

		h.log.Info("Retrieved user orders", zap.Int("count", len(orders)), zap.String("userID", userID.String()))

		return c.Status(200).JSON(orders)
	}
}

// GetUserBalance получает баланс пользователя
func (h *GophermartHandler) GetUserBalance() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		balance, err := h.loyaltyService.GetUserBalance(c.Context(), userID)
		if err != nil {
			return err
		}

		h.log.Info("Retrieved user balance", zap.String("userID", userID.String()), zap.Float64("current", balance.Current), zap.Float64("withdrawn", balance.Withdrawn))

		return c.Status(200).JSON(balance)
	}
}

// WithdrawFunds списывает средства с баланса пользователя
func (h *GophermartHandler) WithdrawFunds() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		var req domain.WithdrawRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Валидация структуры
		if err := h.validator.Validate(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Validation failed: " + err.Error()})
		}

		err := h.loyaltyService.WithdrawFunds(c.Context(), userID, &req)
		if err != nil {
			return err
		}

		h.log.Info("Withdrawal successful", zap.Float64("sum", req.Sum), zap.String("order", req.Order), zap.String("userID", userID.String()))

		return c.Status(200).JSON(fiber.Map{
			"message": "Withdrawal successful",
			"order":   req.Order,
			"sum":     req.Sum,
		})
	}
}

// GetUserWithdrawals получает историю списаний пользователя
func (h *GophermartHandler) GetUserWithdrawals() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		withdrawals, err := h.loyaltyService.GetUserWithdrawals(c.Context(), userID)
		if err != nil {
			return err
		}

		if len(withdrawals) == 0 {
			return c.Status(204).Send(nil)
		}

		h.log.Info("Retrieved user withdrawals", zap.Int("count", len(withdrawals)), zap.String("userID", userID.String()))

		return c.Status(200).JSON(withdrawals)
	}
}
