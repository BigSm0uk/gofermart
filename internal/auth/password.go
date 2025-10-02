package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordManager управляет хешированием паролей
type PasswordManager struct {
	cost int
}

// NewPasswordManager создает новый менеджер паролей
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword хеширует пароль с помощью bcrypt
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword проверяет пароль против хеша
func (pm *PasswordManager) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
