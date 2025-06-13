package raydium

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"deficheck/problem2/internal/types"
	"deficheck/problem2/pkg/solana"
	"deficheck/problem2/pkg/utils"
)

type Client struct {
	solanaClient *solana.Client
}

func NewClient(solanaClient *solana.Client) *Client {
	return &Client{
		solanaClient: solanaClient,
	}
}

const (
	StatusOffset                 = 0
	NonceOffset                  = 8
	MaxOrderOffset               = 16
	DepthOffset                  = 24
	BaseDecimalOffset            = 32
	QuoteDecimalOffset           = 40
	StateOffset                  = 48
	ResetFlagOffset              = 56
	MinSizeOffset                = 64
	VolMaxCutRatioOffset         = 72
	AmountWaveRatioOffset        = 80
	BaseLotSizeOffset            = 88
	QuoteLotSizeOffset           = 96
	MinPriceMultiplierOffset     = 104
	MaxPriceMultiplierOffset     = 112
	SystemDecimalValueOffset     = 120
	MinSeparateNumeratorOffset   = 128
	MinSeparateDenominatorOffset = 136
	TradeFeeNumeratorOffset      = 144
	TradeFeeDenominatorOffset    = 152
	PnlNumeratorOffset           = 160
	PnlDenominatorOffset         = 168
	SwapFeeNumeratorOffset       = 176
	SwapFeeDenominatorOffset     = 184
	BaseNeedTakePnlOffset        = 192
	QuoteNeedTakePnlOffset       = 200
	QuoteTotalPnlOffset          = 208
	BaseTotalPnlOffset           = 216
	PoolOpenTimeOffset           = 232

	// PublicKey fields (32 bytes each)
	BaseVaultOffset     = 336
	QuoteVaultOffset    = 368
	BaseMintOffset      = 400
	QuoteMintOffset     = 432
	LpMintOffset        = 456
	OpenOrdersOffset    = 488
	MarketOffset        = 520
	MarketProgramOffset = 552
	TargetOrdersOffset  = 584
	WithdrawQueueOffset = 616
	LpVaultOffset       = 648
	OwnerOffset         = 680
	PnlOwnerOffset      = 712

	// Compatibility aliases
	CoinVaultOffset = BaseVaultOffset
	PcVaultOffset   = QuoteVaultOffset
	CoinMintOffset  = BaseMintOffset
	PcMintOffset    = QuoteMintOffset
)

func (r *Client) GetPoolInfoOnchain(poolAddress string) (*types.PoolInfo, error) {
	accountInfo, err := r.solanaClient.GetAccountInfo(poolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool account: %w", err)
	}

	value, ok := accountInfo["value"].(map[string]interface{})
	if !ok || value == nil {
		return nil, fmt.Errorf("pool not found")
	}

	dataList, ok := value["data"].([]interface{})
	if !ok || len(dataList) < 2 {
		return nil, fmt.Errorf("invalid account data format")
	}

	base64Data, ok := dataList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid data encoding")
	}

	data, err := solana.DecodeBase64Data(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode account data: %w", err)
	}

	// Validate data length
	if len(data) < 600 {
		return nil, fmt.Errorf("pool data too short: %d bytes (expected at least 600)", len(data))
	}

	baseDecimals := int(binary.LittleEndian.Uint64(data[BaseDecimalOffset : BaseDecimalOffset+8]))
	quoteDecimals := int(binary.LittleEndian.Uint64(data[QuoteDecimalOffset : QuoteDecimalOffset+8]))

	baseMint := extractPubkey(data[CoinMintOffset : CoinMintOffset+32])
	quoteMint := extractPubkey(data[PcMintOffset : PcMintOffset+32])
	baseVault := extractPubkey(data[CoinVaultOffset : CoinVaultOffset+32])
	quoteVault := extractPubkey(data[PcVaultOffset : PcVaultOffset+32])

	baseBalance, err := r.solanaClient.GetTokenAccountBalance(baseVault)
	if err != nil {
		return nil, fmt.Errorf("failed to get base vault balance: %w", err)
	}

	quoteBalance, err := r.solanaClient.GetTokenAccountBalance(quoteVault)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote vault balance: %w", err)
	}

	// Extract amounts from balances
	baseReserve, err := extractTokenAmount(baseBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to extract base reserve: %w", err)
	}

	quoteReserve, err := extractTokenAmount(quoteBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to extract quote reserve: %w", err)
	}

	poolInfo := &types.PoolInfo{
		PoolAddress:   poolAddress,
		BaseToken:     baseMint,
		QuoteToken:    quoteMint,
		BaseReserve:   baseReserve,
		QuoteReserve:  quoteReserve,
		BaseDecimals:  baseDecimals,
		QuoteDecimals: quoteDecimals,
	}

	return poolInfo, nil
}

// GetPoolInfo retrieves pool information from account data (legacy method)
func (r *Client) GetPoolInfo(poolAddress string) (*types.PoolInfo, error) {
	accountInfo, err := r.solanaClient.GetAccountInfo(poolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool account: %w", err)
	}

	value, ok := accountInfo["value"].(map[string]interface{})
	if !ok || value == nil {
		return nil, fmt.Errorf("pool not found")
	}

	dataList, ok := value["data"].([]interface{})
	if !ok || len(dataList) < 2 {
		return nil, fmt.Errorf("invalid account data format")
	}

	base64Data, ok := dataList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid data encoding")
	}

	data, err := solana.DecodeBase64Data(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode account data: %w", err)
	}

	// Validate data length
	if len(data) < 600 {
		return nil, fmt.Errorf("pool data too short: %d bytes (expected at least 600)", len(data))
	}

	baseDecimals := int(binary.LittleEndian.Uint64(data[BaseDecimalOffset : BaseDecimalOffset+8]))
	quoteDecimals := int(binary.LittleEndian.Uint64(data[QuoteDecimalOffset : QuoteDecimalOffset+8]))

	baseMint := extractPubkey(data[CoinMintOffset : CoinMintOffset+32])
	quoteMint := extractPubkey(data[PcMintOffset : PcMintOffset+32])
	baseVault := extractPubkey(data[CoinVaultOffset : CoinVaultOffset+32])
	quoteVault := extractPubkey(data[PcVaultOffset : PcVaultOffset+32])

	vaults, err := r.solanaClient.GetMultipleAccounts([]string{baseVault, quoteVault})
	if err != nil {
		return nil, fmt.Errorf("failed to get vault accounts: %w", err)
	}

	baseReserve, err := getTokenBalance(vaults[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get base reserve: %w", err)
	}

	quoteReserve, err := getTokenBalance(vaults[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get quote reserve: %w", err)
	}

	poolInfo := &types.PoolInfo{
		PoolAddress:   poolAddress,
		BaseToken:     baseMint,
		QuoteToken:    quoteMint,
		BaseReserve:   baseReserve,
		QuoteReserve:  quoteReserve,
		BaseDecimals:  baseDecimals,
		QuoteDecimals: quoteDecimals,
	}

	return poolInfo, nil
}

// CalculatePrice calculates the price for a given quantity and side
func (r *Client) CalculatePrice(pool *types.PoolInfo, quantity *big.Float, side string) (*big.Float, error) {
	// Check for zero reserves
	if pool.BaseReserve.Sign() == 0 || pool.QuoteReserve.Sign() == 0 {
		return nil, fmt.Errorf("pool has zero reserves - pool may be inactive or not initialized")
	}

	baseReserve := new(big.Float).SetInt(pool.BaseReserve)
	quoteReserve := new(big.Float).SetInt(pool.QuoteReserve)

	baseDecimalFactor := new(big.Float).SetFloat64(utils.Pow10(pool.BaseDecimals))
	quoteDecimalFactor := new(big.Float).SetFloat64(utils.Pow10(pool.QuoteDecimals))

	baseReserve.Quo(baseReserve, baseDecimalFactor)
	quoteReserve.Quo(quoteReserve, quoteDecimalFactor)

	if baseReserve.Sign() == 0 || quoteReserve.Sign() == 0 {
		return nil, fmt.Errorf("pool reserves are too small after decimal conversion")
	}

	// AMM: x * y = k
	k := new(big.Float).Mul(baseReserve, quoteReserve)
	var price *big.Float

	if side == "buy" {
		// Buy: spend SOL to get tokens
		newQuoteReserve := new(big.Float).Sub(quoteReserve, quantity)
		newBaseReserve := new(big.Float).Quo(k, newQuoteReserve)
		baseAmount := new(big.Float).Sub(newBaseReserve, baseReserve)
		price = baseAmount
	} else if side == "sell" {
		// Sell: get SOL for tokens
		newQuoteReserve := new(big.Float).Add(quoteReserve, quantity)
		newBaseReserve := new(big.Float).Quo(k, newQuoteReserve)
		baseAmount := new(big.Float).Sub(baseReserve, newBaseReserve)
		price = baseAmount
	} else {
		return nil, fmt.Errorf("invalid side: %s (must be 'buy' or 'sell')", side)
	}

	return price, nil
}

func extractPubkey(data []byte) string {
	return utils.Base58Encode(data)
}

func getTokenBalance(accountData map[string]interface{}) (*big.Int, error) {
	if accountData == nil {
		return big.NewInt(0), nil
	}

	value, ok := accountData["data"].([]interface{})
	if !ok || len(value) < 2 {
		return big.NewInt(0), nil
	}

	base64Data, ok := value[0].(string)
	if !ok {
		return big.NewInt(0), nil
	}

	data, err := solana.DecodeBase64Data(base64Data)
	if err != nil {
		return nil, err
	}

	if len(data) < 8 {
		return big.NewInt(0), nil
	}

	amount := binary.LittleEndian.Uint64(data[0:8])
	result := new(big.Int)
	result.SetUint64(amount)
	return result, nil
}

// extractTokenAmount extracts token amount from getTokenAccountBalance response
func extractTokenAmount(balanceData map[string]interface{}) (*big.Int, error) {
	if balanceData == nil {
		return nil, fmt.Errorf("balance data is nil")
	}

	value, ok := balanceData["value"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid balance response format")
	}

	amountStr, ok := value["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("amount field not found or not a string")
	}

	amount := new(big.Int)
	_, success := amount.SetString(amountStr, 10)
	if !success {
		return nil, fmt.Errorf("failed to parse amount: %s", amountStr)
	}

	return amount, nil
}
