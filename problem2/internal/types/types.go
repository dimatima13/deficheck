package types

import (
	"math/big"
)

type QuoteRequest struct {
	TokenAddress string
	Quantity     *big.Float
	Side         string // "buy" or "sell"
}

type QuoteResponse struct {
	Price          *big.Float
	PriceFormatted string
	TokenSymbol    string
	Decimals       int
	Protocol       string
}

type PoolInfo struct {
	PoolAddress   string
	BaseToken     string
	QuoteToken    string
	BaseReserve   *big.Int
	QuoteReserve  *big.Int
	BaseDecimals  int
	QuoteDecimals int
}

type TokenInfo struct {
	Address  string
	Symbol   string
	Decimals int
}

const (
	ProtocolName = "Raydium"
)
