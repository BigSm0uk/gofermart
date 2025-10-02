package auth

import (
	"testing"
)

func TestPasswordManager(t *testing.T) {
	pm := NewPasswordManager()
	password := "testpassword123"

	// Тест хеширования пароля
	hashedPassword, err := pm.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be equal to original password")
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password should not be empty")
	}

	// Тест проверки правильного пароля
	err = pm.CheckPassword(password, hashedPassword)
	if err != nil {
		t.Errorf("Failed to check correct password: %v", err)
	}

	// Тест проверки неправильного пароля
	wrongPassword := "wrongpassword"
	err = pm.CheckPassword(wrongPassword, hashedPassword)
	if err == nil {
		t.Error("Should fail for wrong password")
	}

	// Тест что каждый раз генерируется разный хеш
	hashedPassword2, err := pm.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password second time: %v", err)
	}

	if hashedPassword == hashedPassword2 {
		t.Error("Each hash should be unique due to salt")
	}

	// Но оба хеша должны валидироваться с правильным паролем
	err = pm.CheckPassword(password, hashedPassword2)
	if err != nil {
		t.Errorf("Second hash should validate with correct password: %v", err)
	}
}
