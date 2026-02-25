package utils

import "testing"

func TestValidateLuhn(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{
			name:     "Valid Luhn number",
			number:   "12345678903",
			expected: true,
		},
		{
			name:     "Valid Luhn number with spaces",
			number:   "1234 5678 903",
			expected: true,
		},
		{
			name:     "Invalid Luhn number",
			number:   "12345678904",
			expected: false,
		},
		{
			name:     "Empty string",
			number:   "",
			expected: false,
		},
		{
			name:     "Non-numeric string",
			number:   "abc123",
			expected: false,
		},
		{
			name:     "Single digit",
			number:   "0",
			expected: true,
		},
		{
			name:     "Another valid Luhn number",
			number:   "4532015112830366",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateLuhn(tt.number)
			if result != tt.expected {
				t.Errorf("ValidateLuhn(%s) = %v, expected %v", tt.number, result, tt.expected)
			}
		})
	}
}
