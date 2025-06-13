package quote

import (
	"fmt"
	"math/big"
	"strings"

	"deficheck/problem2/internal/types"
	"deficheck/problem2/pkg/raydium"
	"deficheck/problem2/pkg/utils"
)

type Service struct {
	raydiumClient *raydium.Client
	raydiumAPI    *raydium.APIClient
	poolCache     map[string]*types.PoolInfo
	useAPI        bool
	useOnchain    bool
}

func NewService(raydiumClient *raydium.Client) *Service {
	return &Service{
		raydiumClient: raydiumClient,
		raydiumAPI:    raydium.NewAPIClient(),
		poolCache:     make(map[string]*types.PoolInfo),
		useAPI:        false,
		useOnchain:    false,
	}
}

func (s *Service) SetUseAPI(useAPI bool) {
	s.useAPI = useAPI
}

func (s *Service) SetUseOnchain(useOnchain bool) {
	s.useOnchain = useOnchain
}

func (s *Service) GetQuote(request *types.QuoteRequest) (*types.QuoteResponse, error) {
	if err := validateRequest(request); err != nil {
		return nil, err
	}

	pool, err := s.findPool(request.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	price, err := s.raydiumClient.CalculatePrice(pool, request.Quantity, request.Side)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate price: %w", err)
	}

	decimals := DetermineDecimals(price)

	priceFormatted := FormatPrice(price, decimals)

	return &types.QuoteResponse{
		Price:          price,
		PriceFormatted: priceFormatted,
		TokenSymbol:    GetTokenSymbol(request.TokenAddress),
		Decimals:       decimals,
		Protocol:       types.ProtocolName,
	}, nil
}

func (s *Service) findPool(tokenAddress string) (*types.PoolInfo, error) {
	if s.useAPI {
		return s.findPoolViaAPI(tokenAddress)
	}

	return s.findPoolHardcoded(tokenAddress)
}

func (s *Service) findPoolViaAPI(tokenAddress string) (*types.PoolInfo, error) {
	fmt.Printf("Searching for pool via Raydium API for token: %s\n", tokenAddress)

	poolV3, err := s.raydiumAPI.FindPoolByTokenV3(tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool via v3 API: %w", err)
	}

	fmt.Printf("Found pool via v3 API: %s (%s-%s) with TVL $%.2f\n",
		poolV3.ID, poolV3.MintA.Symbol, poolV3.MintB.Symbol, poolV3.Tvl)

	// Check cache
	if pool, ok := s.poolCache[poolV3.ID]; ok {
		return pool, nil
	}

	if s.useOnchain {
		fmt.Printf("Using onchain data for pool %s\n", poolV3.ID)
		pool, err := s.raydiumClient.GetPoolInfoOnchain(poolV3.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get onchain pool info: %w", err)
		}
		s.poolCache[poolV3.ID] = pool
		return pool, nil
	}

	pool := s.convertV3PoolToInternal(poolV3)
	s.poolCache[poolV3.ID] = pool
	return pool, nil
}

func (s *Service) findPoolHardcoded(tokenAddress string) (*types.PoolInfo, error) {
	knownPools := map[string]string{
		// USDC-SOL pool (one of the most active)
		"epjfwdd5aufqssqem2qn1xzybapc8g4weggkzwytdt1v": "58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2",
		// USDT-SOL pool
		"es9vmfrzacermjfrf4h2fyd4kconky11mcce8benwnyb": "7XawhbbxtsRcQA8KTkHT9f9nc6d69UwqCDh6U5EEbEmX",
	}

	poolAddress, ok := knownPools[strings.ToLower(tokenAddress)]
	if !ok {
		return nil, fmt.Errorf("no known pool for token %s", tokenAddress)
	}

	fmt.Printf("Using hardcoded pool address: %s for token %s\n", poolAddress, tokenAddress)

	if pool, ok := s.poolCache[poolAddress]; ok {
		return pool, nil
	}

	var pool *types.PoolInfo
	var err error
	if s.useOnchain {
		fmt.Printf("Using onchain data for hardcoded pool %s\n", poolAddress)
		pool, err = s.raydiumClient.GetPoolInfoOnchain(poolAddress)
	} else {
		pool, err = s.raydiumClient.GetPoolInfo(poolAddress)
	}
	if err != nil {
		return nil, err
	}

	s.poolCache[poolAddress] = pool

	return pool, nil
}

func validateRequest(request *types.QuoteRequest) error {
	if request.TokenAddress == "" {
		return fmt.Errorf("token address is required")
	}

	if request.Quantity == nil || request.Quantity.Sign() <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	side := strings.ToLower(request.Side)
	if side != "buy" && side != "sell" {
		return fmt.Errorf("side must be 'buy' or 'sell'")
	}

	return nil
}

func DetermineDecimals(price *big.Float) int {
	priceFloat, _ := price.Float64()

	if priceFloat < 0.01 {
		return 8
	} else if priceFloat < 1 {
		return 6
	} else if priceFloat < 100 {
		return 4
	} else {
		return 2
	}
}

func FormatPrice(price *big.Float, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	priceFloat, _ := price.Float64()
	return fmt.Sprintf(format, priceFloat)
}

func GetTokenSymbol(address string) string {
	knownTokens := map[string]string{
		"epjfwdd5aufqssqem2qn1xzybapc8g4weggkzwytdt1v": "USDC",
		"es9vmfrzacermjfrf4h2fyd4kconky11mcce8benwnyb": "USDT",
		"so11111111111111111111111111111111111111112":  "SOL",
	}

	if symbol, ok := knownTokens[strings.ToLower(address)]; ok {
		return symbol
	}

	if len(address) > 8 {
		return address[:4] + "..." + address[len(address)-4:]
	}
	return address
}

func (s *Service) convertV3PoolToInternal(poolV3 *raydium.PoolInfoData) *types.PoolInfo {
	baseReserve := new(big.Int)
	quoteReserve := new(big.Int)

	baseMultiplier := new(big.Float).SetFloat64(utils.Pow10(poolV3.MintA.Decimals))
	quoteMultiplier := new(big.Float).SetFloat64(utils.Pow10(poolV3.MintB.Decimals))

	baseAmount := new(big.Float).SetFloat64(poolV3.MintAmountA)
	quoteAmount := new(big.Float).SetFloat64(poolV3.MintAmountB)

	baseAmount.Mul(baseAmount, baseMultiplier)
	quoteAmount.Mul(quoteAmount, quoteMultiplier)

	baseAmount.Int(baseReserve)
	quoteAmount.Int(quoteReserve)

	return &types.PoolInfo{
		PoolAddress:   poolV3.ID,
		BaseToken:     poolV3.MintA.Address,
		QuoteToken:    poolV3.MintB.Address,
		BaseReserve:   baseReserve,
		QuoteReserve:  quoteReserve,
		BaseDecimals:  poolV3.MintA.Decimals,
		QuoteDecimals: poolV3.MintB.Decimals,
	}
}
