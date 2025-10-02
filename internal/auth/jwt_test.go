package auth

import (
	"testing"
	"time"
)

func TestJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := time.Hour
	jwtManager := NewJWTManager(secretKey, tokenDuration)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	login := "testuser"

	// Тест генерации токена
	token, err := jwtManager.GenerateToken(userID, login)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}

	// Тест валидации токена
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, claims.UserID)
	}

	if claims.Login != login {
		t.Errorf("Expected login %s, got %s", login, claims.Login)
	}

	// Тест валидации неверного токена
	wrongToken := "invalid.token.here"
	_, err = jwtManager.ValidateToken(wrongToken)
	if err == nil {
		t.Error("Should fail for invalid token")
	}

	// Тест валидации токена с неверным секретом
	wrongSecretManager := NewJWTManager("wrong-secret", tokenDuration)
	wrongToken, _ = wrongSecretManager.GenerateToken(userID, login)
	_, err = jwtManager.ValidateToken(wrongToken)
	if err == nil {
		t.Error("Should fail for token with wrong secret")
	}

	// Тест истечения токена
	shortDurationManager := NewJWTManager(secretKey, time.Nanosecond)
	expiredToken, err := shortDurationManager.GenerateToken(userID, login)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	time.Sleep(time.Millisecond) // Ждем истечения токена

	_, err = shortDurationManager.ValidateToken(expiredToken)
	if err == nil {
		t.Error("Should fail for expired token")
	}
}
