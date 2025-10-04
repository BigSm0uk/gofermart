package service

import (
	"context"
	"fmt"

	"github.com/BigSm0uk/gofermart/internal/auth"
	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/google/uuid"
)

// UserService представляет сервис для работы с пользователями
type UserService struct {
	userRepo        *repo.UserRepository
	passwordManager *auth.PasswordManager
	jwtManager      *auth.JWTManager
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo *repo.UserRepository, passwordManager *auth.PasswordManager, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo:        userRepo,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
	}
}

// RegisterUser регистрирует нового пользователя
func (s *UserService) RegisterUser(ctx context.Context, req *domain.UserCreateRequest) (*domain.User, string, error) {
	// Хешируем пароль
	hashedPassword, err := s.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user, err := s.userRepo.CreateUser(ctx, req.Login, hashedPassword)
	if err != nil {
		return nil, "", err
	}

	// Генерируем JWT токен
	token, err := s.jwtManager.GenerateToken(user.ID.String(), user.Login)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// LoginUser аутентифицирует пользователя
func (s *UserService) LoginUser(ctx context.Context, req *domain.UserLoginRequest) (*domain.User, string, error) {
	// Получаем пользователя по логину
	user, err := s.userRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, "", domain.ErrInvalidCredentials
		}
		return nil, "", err
	}

	// Проверяем пароль
	err = s.passwordManager.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	// Генерируем JWT токен
	token, err := s.jwtManager.GenerateToken(user.ID.String(), user.Login)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// GetUserByID получает пользователя по ID
func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}
