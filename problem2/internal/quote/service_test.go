package quote

import (
	"math/big"
	"testing"

	"deficheck/problem2/internal/types"
	"deficheck/problem2/pkg/raydium"
	"deficheck/problem2/pkg/solana"
)

func TestServiceWithV3API(t *testing.T) {
	// Create real clients
	solanaClient := solana.NewClient("https://api.mainnet-beta.solana.com")
	raydiumClient := raydium.NewClient(solanaClient)
	service := NewService(raydiumClient)

	// Enable API usage
	service.SetUseAPI(true)

	t.Run("Get quote for USDC using v3 API", func(t *testing.T) {
		request := &types.QuoteRequest{
			TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
			Quantity:     big.NewFloat(1.0),                              // 1 SOL
			Side:         "buy",
		}

		response, err := service.GetQuote(request)
		if err != nil {
			t.Logf("Could not get quote (API might be unavailable): %v", err)
			t.Skip("Skipping test - API might be unavailable")
		}

		// Check response
		if response == nil {
			t.Fatal("Expected response, got nil")
		}

		if response.Price == nil || response.Price.Sign() <= 0 {
			t.Error("Expected positive price")
		}

		t.Logf("Quote for 1 SOL -> USDC: %s USDC", response.PriceFormatted)
		t.Logf("Token: %s", response.TokenSymbol)
		t.Logf("Protocol: %s", response.Protocol)
	})

}

func TestV3APIPoolSelection(t *testing.T) {
	apiClient := raydium.NewAPIClient()

	t.Run("Verify SOL pool selection logic", func(t *testing.T) {
		// Search for USDC pool
		pool, err := apiClient.FindPoolByTokenV3("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
		if err != nil {
			t.Skip("API unavailable")
		}

		solMint := "So11111111111111111111111111111111111111112"
		
		// Verify it's paired with SOL
		isPairedWithSOL := pool.MintA.Address == solMint || pool.MintB.Address == solMint
		
		if !isPairedWithSOL {
			t.Error("Expected pool to be paired with SOL")
		}
		
		t.Logf("Found SOL pool: %s", pool.ID)
		t.Logf("  Type: %s", pool.Type)
		t.Logf("  Pair: %s-%s", pool.MintA.Symbol, pool.MintB.Symbol)
		t.Logf("  TVL: $%.2f", pool.Tvl)
		t.Logf("  24h Volume: $%.2f", pool.Day.Volume)
	})
}