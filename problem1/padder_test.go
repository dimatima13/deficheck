package problem1

import (
	"testing"
)

func TestPadNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{
			name:     "James Bond example",
			input:    "James Bond 7",
			width:    3,
			expected: "James Bond 007",
		},
		{
			name:     "PI decimal example",
			input:    "PI=3.14",
			width:    2,
			expected: "PI=03.14",
		},
		{
			name:     "Time example with single digit",
			input:    "It's 3:13pm",
			width:    2,
			expected: "It's 03:13pm",
		},
		{
			name:     "Time example with two digits",
			input:    "It's 12:13pm",
			width:    2,
			expected: "It's 12:13pm",
		},
		{
			name:     "Mixed alphanumeric",
			input:    "99UR1337",
			width:    6,
			expected: "000099UR001337",
		},
		{
			name:     "Multiple numbers",
			input:    "1 2 3 4 5",
			width:    3,
			expected: "001 002 003 004 005",
		},
		{
			name:     "Numbers with leading zeros",
			input:    "007 is agent 007",
			width:    4,
			expected: "0007 is agent 0007",
		},
		{
			name:     "Empty string",
			input:    "",
			width:    3,
			expected: "",
		},
		{
			name:     "No numbers",
			input:    "Hello World!",
			width:    3,
			expected: "Hello World!",
		},
		{
			name:     "Decimal numbers",
			input:    "Price is 10.99 dollars",
			width:    3,
			expected: "Price is 010.99 dollars",
		},
		{
			name:     "Numbers at start and end",
			input:    "100 bottles of beer 99",
			width:    4,
			expected: "0100 bottles of beer 0099",
		},
		{
			name:     "Width smaller than number",
			input:    "Year 2024",
			width:    2,
			expected: "Year 2024",
		},
		{
			name:     "Complex decimal",
			input:    "Values: 1.23 and 45.6789",
			width:    3,
			expected: "Values: 001.23 and 045.6789",
		},
		{
			name:     "Number with punctuation",
			input:    "Call 911!",
			width:    4,
			expected: "Call 0911!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadNumbers(tt.input, tt.width)
			if result != tt.expected {
				t.Errorf("PadNumbers(%q, %d) = %q; want %q",
					tt.input, tt.width, result, tt.expected)
			}
		})
	}
}

func TestPadNumbersEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{
			name:     "Zero width",
			input:    "Test 123",
			width:    0,
			expected: "Test 123",
		},
		{
			name:     "Negative width",
			input:    "Test 123",
			width:    -1,
			expected: "Test 123",
		},
		{
			name:     "Very large width",
			input:    "Number 5",
			width:    10,
			expected: "Number 0000000005",
		},
		{
			name:     "Unicode characters",
			input:    "Цена: 100 рублей",
			width:    4,
			expected: "Цена: 0100 рублей",
		},
		{
			name:     "Special characters between numbers",
			input:    "12-34-56",
			width:    3,
			expected: "012-034-056",
		},
		{
			name:     "Numbers in parentheses",
			input:    "(123) 456-7890",
			width:    4,
			expected: "(0123) 0456-7890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadNumbers(tt.input, tt.width)
			if result != tt.expected {
				t.Errorf("PadNumbers(%q, %d) = %q; want %q",
					tt.input, tt.width, result, tt.expected)
			}
		})
	}
}

func BenchmarkPadNumbers(b *testing.B) {
	testCases := []struct {
		name  string
		input string
		width int
	}{
		{"Short", "Test 123", 3},
		{"Medium", "The year 2024 has 365 days and 12 months", 4},
		{"Long", "Numbers: 1 22 333 4444 55555 666666 7777777 88888888 999999999", 10},
		{"ManySmall", "1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20", 3},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = PadNumbers(tc.input, tc.width)
			}
		})
	}
}