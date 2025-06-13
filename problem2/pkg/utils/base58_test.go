package utils

import (
	"bytes"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Empty input",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "Single zero byte",
			input:    []byte{0},
			expected: "1",
		},
		{
			name:     "Multiple zero bytes",
			input:    []byte{0, 0, 0},
			expected: "111",
		},
		{
			name:     "Simple bytes",
			input:    []byte{1, 2, 3},
			expected: "Ldp",
		},
		{
			name:     "Solana public key example",
			input:    []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: "CiDwVBFgWV9E5MvXWoLgnEgn2hK7rJikbvfWavzAQz3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base58Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Base58Encode(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBase58Decode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: []byte{},
			wantErr:  false,
		},
		{
			name:     "Single 1 (zero byte)",
			input:    "1",
			expected: []byte{0},
			wantErr:  false,
		},
		{
			name:     "Multiple 1s (zero bytes)",
			input:    "111",
			expected: []byte{0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "Simple string",
			input:    "Ldp",
			expected: []byte{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "Invalid character",
			input:    "0OIl",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Solana public key",
			input:    "CiDwVBFgWV9E5MvXWoLgnEgn2hK7rJikbvfWavzAQz3",
			expected: []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Base58Decode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Base58Decode(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !bytes.Equal(result, tt.expected) {
				t.Errorf("Base58Decode(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBase58RoundTrip(t *testing.T) {
	// Test that encoding and then decoding gives back the original
	testCases := [][]byte{
		{},
		{0},
		{0, 0, 0},
		{255},
		{1, 2, 3, 4, 5},
		{0, 1, 2, 3, 4, 5, 0},
		// Typical Solana public key (32 bytes)
		{
			1, 2, 3, 4, 5, 6, 7, 8,
			9, 10, 11, 12, 13, 14, 15, 16,
			17, 18, 19, 20, 21, 22, 23, 24,
			25, 26, 27, 28, 29, 30, 31, 32,
		},
	}

	for i, original := range testCases {
		encoded := Base58Encode(original)
		decoded, err := Base58Decode(encoded)
		if err != nil {
			t.Errorf("Test case %d: Base58Decode failed: %v", i, err)
			continue
		}
		if !bytes.Equal(decoded, original) {
			t.Errorf("Test case %d: Round trip failed. Original: %v, Decoded: %v", i, original, decoded)
		}
	}
}