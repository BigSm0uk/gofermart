package repo

import (
	"context"
	"fmt"

	"github.com/BigSm0uk/gofermart/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository представляет репозиторий для работы с пользователями
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, login, passwordHash string) (*domain.User, error) {
	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id, login, password_hash, created_at, updated_at`

	var user domain.User
	err := r.db.QueryRow(ctx, query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// Проверяем на дублирование логина
		if isUniqueViolation(err) {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetUserByLogin получает пользователя по логину
func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `
		SELECT id, login, password_hash, created_at, updated_at
		FROM users
		WHERE login = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}

	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, login, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// isUniqueViolation проверяет, является ли ошибка нарушением уникальности
func isUniqueViolation(err error) bool {
	// В pgx v5 можно использовать более точную проверку
	// Пока используем простое сравнение строк
	return err != nil && (fmt.Sprintf("%v", err) == "duplicate key value violates unique constraint" ||
		fmt.Sprintf("%v", err) == "duplicate key value violates unique constraint \"users_login_key\"")
}
