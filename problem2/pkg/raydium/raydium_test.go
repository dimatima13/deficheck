package raydium

import (
	"math/big"
	"testing"
	
	"deficheck/problem2/internal/types"
)

func TestCalculatePrice(t *testing.T) {
	// Create a mock pool with known reserves
	pool := &types.PoolInfo{
		PoolAddress:   "mockpool",
		BaseToken:     "basetoken",
		QuoteToken:    "quotetoken",
		BaseReserve:   big.NewInt(1000000000000), // 1000 with 9 decimals
		QuoteReserve:  big.NewInt(50000000000),   // 50 with 9 decimals
		BaseDecimals:  9,
		QuoteDecimals: 9,
	}

	client := &Client{}

	tests := []struct {
		name        string
		pool        *types.PoolInfo
		quantity    *big.Float
		side        string
		wantErr     bool
		priceCheck  func(*big.Float) bool
	}{
		{
			name:     "Buy small amount",
			pool:     pool,
			quantity: big.NewFloat(1),
			side:     "buy",
			wantErr:  false,
			priceCheck: func(price *big.Float) bool {
				// For buying 1 quote token from a 50/1000 pool
				// We need to spend some base tokens (SOL)
				// With reserves 50 base and 1000 quote, price should be around 20 SOL
				// (due to AMM formula impact on small pool)
				return price.Sign() > 0 && price.Cmp(big.NewFloat(10)) > 0 && price.Cmp(big.NewFloat(30)) < 0
			},
		},
		{
			name:     "Sell small amount",
			pool:     pool,
			quantity: big.NewFloat(1),
			side:     "sell",
			wantErr:  false,
			priceCheck: func(price *big.Float) bool {
				// For selling 1 quote token to a 50/1000 pool
				// We receive some base tokens (SOL)
				// With reserves 50 base and 1000 quote, price should be around 19.6 SOL
				return price.Sign() > 0 && price.Cmp(big.NewFloat(10)) > 0 && price.Cmp(big.NewFloat(30)) < 0
			},
		},
		{
			name:     "Invalid side",
			pool:     pool,
			quantity: big.NewFloat(1),
			side:     "invalid",
			wantErr:  true,
			priceCheck: func(price *big.Float) bool {
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := client.CalculatePrice(tt.pool, tt.quantity, tt.side)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculatePrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !tt.priceCheck(price) {
				t.Errorf("CalculatePrice() price = %v, failed price check", price)
			}
		})
	}
}

func TestExtractPubkey(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Simple bytes",
			input:    []byte{1, 2, 3, 4},
			expected: "2VfUX",
		},
		{
			name:     "All zeros (32 bytes)",
			input:    make([]byte, 32),
			expected: "11111111111111111111111111111111",
		},
		{
			name:     "Typical Solana address",
			input:    []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: "CiDwVBFgWV9E5MvXWoLgnEgn2hK7rJikbvfWavzAQz3",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPubkey(tt.input)
			if result != tt.expected {
				t.Errorf("extractPubkey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetTokenBalance(t *testing.T) {
	tests := []struct {
		name        string
		accountData map[string]interface{}
		expected    int64
		wantErr     bool
	}{
		{
			name:        "Nil account",
			accountData: nil,
			expected:    0,
			wantErr:     false,
		},
		{
			name: "Valid token account",
			accountData: map[string]interface{}{
				"data": []interface{}{
					"AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==", // 1 encoded as uint64 little endian in base64
					"base64",
				},
			},
			expected: 1,
			wantErr:  false,
		},
		{
			name:        "Empty data",
			accountData: map[string]interface{}{},
			expected:    0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getTokenBalance(tt.accountData)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result.Int64() != tt.expected {
				t.Errorf("getTokenBalance() = %v, want %v", result.Int64(), tt.expected)
			}
		})
	}
}