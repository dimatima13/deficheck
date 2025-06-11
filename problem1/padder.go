package problem1

import (
	"fmt"
	"strconv"
	"strings"
)

// PadNumbers takes a string and an integer X, returns a string with whole numbers
// left-padded with zeros to X characters
func PadNumbers(input string, width int) string {
	var result strings.Builder
	i := 0

	for i < len(input) {
		if i < len(input) && input[i] >= '0' && input[i] <= '9' {
			start := i
			for i < len(input) && input[i] >= '0' && input[i] <= '9' {
				i++
			}
			numStr := input[start:i]

			// Check for decimal fraction (dot before number means it's a fractional part)
			isDecimalPart := start > 0 && input[start-1] == '.'

			if !isDecimalPart {
				// This is a whole number (possibly with fractional part after it)
				num, _ := strconv.Atoi(numStr)
				result.WriteString(fmt.Sprintf("%0*d", width, num))
			} else {
				// This is the fractional part of a decimal number - leave as is
				result.WriteString(numStr)
			}
		} else {
			// Copy character as is
			result.WriteByte(input[i])
			i++
		}
	}

	return result.String()
}
