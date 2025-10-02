package luhn

import (
	"strconv"
	"strings"
)

// Validate проверяет номер заказа по алгоритму Луна
func Validate(number string) bool {
	// Удаляем все пробелы и проверяем, что строка не пустая
	number = strings.ReplaceAll(number, " ", "")
	if len(number) <= 1 {
		return false
	}

	// Проверяем, что все символы - цифры
	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Алгоритм Луна
	sum := 0
	isEven := false

	// Итерируемся по цифрам справа налево
	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))

		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEven = !isEven
	}

	// Сумма должна быть кратна 10
	return sum%10 == 0
}
