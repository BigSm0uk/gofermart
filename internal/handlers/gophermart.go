package handlers

import (
	"log"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// GophermartHandler представляет HTTP хендлеры для gophermart
type GophermartHandler struct {
	userService    *service.UserService
	orderService   *service.OrderService
	loyaltyService *service.LoyaltyService
}

// NewGophermartHandler создает новый хендлер
func NewGophermartHandler(userService *service.UserService, orderService *service.OrderService, loyaltyService *service.LoyaltyService) *GophermartHandler {
	return &GophermartHandler{
		userService:    userService,
		orderService:   orderService,
		loyaltyService: loyaltyService,
	}
}

// RegisterUser регистрирует нового пользователя
func (h *GophermartHandler) RegisterUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.UserCreateRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Валидация
		if req.Login == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Login and password are required"})
		}

		if len(req.Password) < 6 {
			return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters"})
		}

		user, token, err := h.userService.RegisterUser(c.Context(), &req)
		if err != nil {
			return err
		}

		log.Printf("User %s registered successfully", user.Login)

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
func (h *GophermartHandler) LoginUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.UserLoginRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Валидация
		if req.Login == "" || req.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Login and password are required"})
		}

		user, token, err := h.userService.LoginUser(c.Context(), &req)
		if err != nil {
			return err
		}

		log.Printf("User %s logged in successfully", user.Login)

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
func (h *GophermartHandler) CreateOrder() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)
		orderNumber := string(c.Body())

		if len(orderNumber) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Order number is required"})
		}

		order, err := h.orderService.CreateOrder(c.Context(), orderNumber, userID)
		if err != nil {
			return err
		}

		// Определяем статус ответа
		status := 202 // Новый заказ
		if order.Status != domain.OrderStatusNew {
			status = 200 // Заказ уже существовал
		}

		log.Printf("Order %s created for user %s", order.Number, userID)

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
func (h *GophermartHandler) GetUserOrders() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		orders, err := h.orderService.GetUserOrders(c.Context(), userID)
		if err != nil {
			return err
		}

		if len(orders) == 0 {
			return c.Status(204).Send(nil)
		}

		log.Printf("Retrieved %d orders for user %s", len(orders), userID)

		return c.Status(200).JSON(orders)
	}
}

// GetUserBalance получает баланс пользователя
func (h *GophermartHandler) GetUserBalance() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		balance, err := h.loyaltyService.GetUserBalance(c.Context(), userID)
		if err != nil {
			return err
		}

		log.Printf("Retrieved balance for user %s: current=%.2f, withdrawn=%.2f", userID, balance.Current, balance.Withdrawn)

		return c.Status(200).JSON(balance)
	}
}

// WithdrawFunds списывает средства с баланса пользователя
func (h *GophermartHandler) WithdrawFunds() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		var req domain.WithdrawRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Валидация
		if req.Order == "" || req.Sum <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Order number and positive sum are required"})
		}

		err := h.loyaltyService.WithdrawFunds(c.Context(), userID, &req)
		if err != nil {
			return err
		}

		log.Printf("Withdrew %.2f for order %s from user %s", req.Sum, req.Order, userID)

		return c.Status(200).JSON(fiber.Map{
			"message": "Withdrawal successful",
			"order":   req.Order,
			"sum":     req.Sum,
		})
	}
}

// GetUserWithdrawals получает историю списаний пользователя
func (h *GophermartHandler) GetUserWithdrawals() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)

		withdrawals, err := h.loyaltyService.GetUserWithdrawals(c.Context(), userID)
		if err != nil {
			return err
		}

		if len(withdrawals) == 0 {
			return c.Status(204).Send(nil)
		}

		log.Printf("Retrieved %d withdrawals for user %s", len(withdrawals), userID)

		return c.Status(200).JSON(withdrawals)
	}
}
