This repository contains solutions for programming challenges implemented in Go.

## Problem 1: Number Padding

This solution takes an input string and pads all whole numbers found in the string with leading zeros to a specified width.

### Implementation Details

The solution uses a character-by-character parsing approach to identify whole numbers while correctly handling decimal numbers. Key features:

- **Decimal handling**: Numbers that are part of decimal values (e.g., "3.14") are treated specially - only the integer part is padded
- **Performance**: O(n) time complexity where n is the length of the input string
- **Memory**: O(n) space complexity for the output string
- **Unicode support**: Works correctly with Unicode characters

### Running the Solution

#### Run the demo program:
```bash
cd problem1/cmd
go run main.go "James Bond 7" 3
```

#### Run tests:
```bash
cd problem1
go test -v
```

#### Run benchmarks:
```bash
cd problem1
go test -bench=.
```

### Example Usage

```go
import "deficheck/problem1"

result := problem1.PadNumbers("James Bond 7", 3)
// Output: "James Bond 007"
```

### Performance Analysis

**Time Complexity**: O(n) where n is the length of the input string
- Single pass through the string
- Each character is processed exactly once

**Space Complexity**: O(n) for the output string
- Uses strings.Builder for efficient string concatenation
- No additional data structures required

**Optimization Decisions**:
- Chose time efficiency over space by using strings.Builder
- Character-by-character parsing avoids regex overhead for simple cases
- No preprocessing or multiple passes required

### Dependencies

No external dependencies. Uses only Go standard library:
- `fmt`: For formatting padded numbers
- `strconv`: For string to integer conversion
- `strings`: For efficient string building
- `regexp`: For alternative implementation (not used in main solution)