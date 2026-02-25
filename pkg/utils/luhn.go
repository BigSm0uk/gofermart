package utils

import (
	"strconv"
	"strings"
)

// ValidateLuhn проверяет номер заказа по алгоритму Луна
func ValidateLuhn(number string) bool {
	// Убираем все пробелы
	number = strings.ReplaceAll(number, " ", "")
	
	// Проверяем, что строка содержит только цифры
	if _, err := strconv.Atoi(number); err != nil {
		return false
	}
	
	// Проверяем длину (минимум 1 цифра)
	if len(number) < 1 {
		return false
	}
	
	// Алгоритм Луна
	sum := 0
	alternate := false
	
	// Проходим по цифрам справа налево
	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))
		
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}
		
		sum += digit
		alternate = !alternate
	}
	
	return sum%10 == 0
}
