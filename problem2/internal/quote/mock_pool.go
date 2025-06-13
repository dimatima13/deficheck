package quote

import (
	"fmt"
	"math/big"
	"strings"
	
	"deficheck/problem2/internal/types"
)

// GetMockPool returns a mock pool with realistic data for testing
func GetMockPool(tokenAddress string) (*types.PoolInfo, error) {
	switch strings.ToLower(tokenAddress) {
	case "epjfwdd5aufqssqem2qn1xzybapc8g4weggkzwytdt1v":
		// USDC-SOL mock pool
		return &types.PoolInfo{
			PoolAddress:   "6UmmUiYoBjSrhakAobJw8BvkmJtDVxaeBtbt7rxWo1mg",
			BaseToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
			QuoteToken:    "So11111111111111111111111111111111111111112",    // SOL
			BaseReserve:   big.NewInt(50000000000000), // 50,000 USDC (6 decimals)
			QuoteReserve:  big.NewInt(1000000000000),  // 1,000 SOL (9 decimals)
			BaseDecimals:  6,
			QuoteDecimals: 9,
		}, nil
	case "es9vmfrzacermjfrf4h2fyd4kconky11mcce8benwnyb":
		// USDT-SOL mock pool
		return &types.PoolInfo{
			PoolAddress:   "7XawhbbxtsRcQA8KTkHT9f9nc6d69UwqCDh6U5EEbEmX",
			BaseToken:     "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", // USDT
			QuoteToken:    "So11111111111111111111111111111111111111112",    // SOL
			BaseReserve:   big.NewInt(30000000000000), // 30,000 USDT (6 decimals)
			QuoteReserve:  big.NewInt(600000000000),   // 600 SOL (9 decimals)
			BaseDecimals:  6,
			QuoteDecimals: 9,
		}, nil
	default:
		return nil, fmt.Errorf("no mock pool available for token %s", tokenAddress)
	}
}