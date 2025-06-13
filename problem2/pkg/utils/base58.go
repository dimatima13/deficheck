package utils

import (
	"fmt"
	"math/big"
)

// Base58 alphabet used by Solana
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var base58AlphabetIdx = map[byte]int{}

func init() {
	for i := 0; i < len(base58Alphabet); i++ {
		base58AlphabetIdx[base58Alphabet[i]] = i
	}
}

// Base58Encode encodes bytes to base58 string
func Base58Encode(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	// Count leading zeros
	leadingZeros := 0
	for _, b := range input {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	// Convert bytes to big integer
	bigInt := new(big.Int).SetBytes(input)

	// Convert to base58
	result := make([]byte, 0, len(input)*2)
	base := big.NewInt(58)
	zero := big.NewInt(0)

	for bigInt.Cmp(zero) > 0 {
		mod := new(big.Int)
		bigInt.DivMod(bigInt, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}

	// Add leading zeros (represented as '1' in base58)
	for i := 0; i < leadingZeros; i++ {
		result = append(result, '1')
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// Base58Decode decodes base58 string to bytes
func Base58Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return []byte{}, nil
	}

	// Count leading ones (zeros in the result)
	leadingOnes := 0
	for _, c := range input {
		if c == '1' {
			leadingOnes++
		} else {
			break
		}
	}

	// Decode
	bigInt := big.NewInt(0)
	base := big.NewInt(58)

	for _, c := range input {
		idx, ok := base58AlphabetIdx[byte(c)]
		if !ok {
			return nil, fmt.Errorf("invalid base58 character: %c", c)
		}

		bigInt.Mul(bigInt, base)
		bigInt.Add(bigInt, big.NewInt(int64(idx)))
	}

	// Convert to bytes
	result := bigInt.Bytes()

	// Add leading zeros
	for i := 0; i < leadingOnes; i++ {
		result = append([]byte{0}, result...)
	}

	return result, nil
}
