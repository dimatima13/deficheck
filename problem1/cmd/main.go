package main

import (
	"fmt"
	"os"
	"strconv"

	"deficheck/problem1"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go \"<input string>\" <width>")
		fmt.Println("Example: go run main.go \"James Bond 7\" 3")
		os.Exit(1)
	}

	input := os.Args[1]
	width, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: Invalid width parameter: %s\n", os.Args[2])
		os.Exit(1)
	}

	result := problem1.PadNumbers(input, width)
	fmt.Printf("Input:  %s\n", input)
	fmt.Printf("Width:  %d\n", width)
	fmt.Printf("Output: %s\n", result)
}
