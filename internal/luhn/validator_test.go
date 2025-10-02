package luhn

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{
			name:     "Valid credit card number",
			number:   "4532015112830366",
			expected: true,
		},
		{
			name:     "Valid order number from specification",
			number:   "9278923470",
			expected: true,
		},
		{
			name:     "Valid order number from specification",
			number:   "12345678903",
			expected: true,
		},
		{
			name:     "Valid order number from specification",
			number:   "2377225624",
			expected: true,
		},
		{
			name:     "Invalid number",
			number:   "1234567890",
			expected: false,
		},
		{
			name:     "Invalid number",
			number:   "1111111111",
			expected: false,
		},
		{
			name:     "Empty string",
			number:   "",
			expected: false,
		},
		{
			name:     "Single digit",
			number:   "5",
			expected: false,
		},
		{
			name:     "Contains non-digits",
			number:   "1234abc5678",
			expected: false,
		},
		{
			name:     "Contains spaces",
			number:   "4532 0151 1283 0366",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.number)
			if result != tt.expected {
				t.Errorf("Validate(%s) = %v, expected %v", tt.number, result, tt.expected)
			}
		})
	}
}
