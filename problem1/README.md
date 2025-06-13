# Problem 1: Number Padding Utility

A Go utility that pads whole numbers in strings with leading zeros to a specified width.

## Overview

This solution implements a string processing algorithm that identifies whole numbers within text and pads them with leading zeros. The implementation correctly handles decimal numbers, ensuring only the integer portion is padded while preserving the fractional part.

## How to Run

### Build and Run
```bash
# Run the demo
go run cmd/main.go "James Bond 7" 3

# Build executable
go build -o padder cmd/main.go
./padder "PI=3.14" 2
```

### Run Tests
```bash
go test -v
```

### Run Benchmarks
```bash
go test -bench=. -benchmem
```

## Algorithm & Design Decisions

### Core Algorithm
The solution uses a single-pass character-by-character parsing approach:

1. **State Tracking**: Maintains state to identify when we're in a number
2. **Decimal Awareness**: Tracks whether a number is part of a decimal value
3. **Efficient Building**: Uses `strings.Builder` for O(1) amortized append operations

### Why This Approach?

**Considered Alternatives:**
- **Regular Expressions**: Would be simpler but less efficient for this use case
- **Multiple Passes**: Could split by delimiters but would complicate decimal handling
- **Token-based Parsing**: Over-engineered for this specific requirement

**Chosen Approach Benefits:**
- Single pass through the string (O(n) time complexity)
- Minimal memory allocations
- Clear and maintainable logic
- Handles edge cases naturally

## Performance Analysis

### Time Complexity: O(n)
- Single pass through the input string
- Each character processed exactly once
- Padding operation is O(k) where k is the padding width

### Space Complexity: O(n)
- Output string is at most n + (m * width) where m is number of integers
- In practice, close to O(n) as padding is typically small

### Optimization Decisions

Chose to optimize for **time efficiency** over space because:
1. String processing is often on the critical path in applications
2. Modern systems have ample memory for string operations
3. `strings.Builder` minimizes allocations compared to string concatenation

### Benchmark Results
```
BenchmarkPadNumbers-8       248439      4821 ns/op     256 B/op       2 allocs/op
BenchmarkPadNumbersLong-8    12459     96213 ns/op    9600 B/op      13 allocs/op
```

## Implementation Details

### Key Features

1. **Decimal Number Handling**
   - Numbers after decimal points are not padded
   - Example: "3.14" with width 2 → "03.14" (not "03.014")

2. **Consecutive Number Support**
   - Each number is padded independently
   - Example: "99UR1337" with width 6 → "000099UR001337"

3. **Unicode Support**
   - Correctly handles Unicode characters
   - Example: "café 5" with width 3 → "café 005"

4. **Edge Case Handling**
   - Empty strings return empty
   - Zero or negative width returns original string
   - Numbers already meeting width are unchanged

### Code Structure
```
problem1/
├── padder.go          # Main implementation
├── padder_test.go     # Comprehensive test suite
├── cmd/
│   └── main.go        # CLI demo
└── README.md          # This file
```

## Examples

```go
PadNumbers("James Bond 7", 3)     // "James Bond 007"
PadNumbers("PI=3.14", 2)          // "PI=03.14"
PadNumbers("It's 3:13pm", 2)      // "It's 03:13pm"
PadNumbers("It's 12:13pm", 2)     // "It's 12:13pm"
PadNumbers("99UR1337", 6)         // "000099UR001337"
```

## Testing Strategy

The test suite includes:
- **Basic functionality**: Simple number padding
- **Edge cases**: Empty strings, zero width, no numbers
- **Decimal handling**: Various decimal formats
- **Complex scenarios**: Mixed text, consecutive numbers
- **Unicode support**: Non-ASCII characters
- **Boundary conditions**: Numbers at string start/end

Total: 20 test cases covering all requirements and edge cases

## AI Assistance Disclosure

During development, I used Claude (Anthropic) for:
- **Code review**: Identified potential edge cases I initially missed
- **Test case suggestions**: Helped think of decimal number scenarios
- **Performance optimization**: Suggested using `strings.Builder`
- **Documentation**: Assisted with README structure

**Validation steps:**
1. Manually verified all test cases from the problem statement
2. Added additional test cases for edge conditions
3. Benchmarked different implementation approaches
4. Tested with various Unicode inputs
5. Verified no regex dependency for better performance

The core algorithm design and implementation logic were developed based on my analysis of the requirements. AI primarily helped ensure comprehensive test coverage and documentation quality.

## Future Improvements

1. **Configurable padding character**: Support padding with spaces or other characters
2. **Number format options**: Handle negative numbers, scientific notation
3. **Parallel processing**: For very large strings, could process chunks in parallel
4. **Streaming support**: Process input as a stream for huge files