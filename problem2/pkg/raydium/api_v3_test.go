package raydium

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetPoolInfoV3(t *testing.T) {
	client := NewAPIClient()

	// Test with known pool ID (USDC-SOL)
	t.Run("Get USDC-SOL pool info", func(t *testing.T) {
		poolID := "6UmmUiYoBjSrhakAobJw8BvkmJtDVxaeBtbt7rxWo1mg"

		poolInfo, err := client.GetPoolInfoV3(poolID)
		if err != nil {
			t.Logf("Could not get pool info from v3 API (this is normal if API is unavailable): %v", err)
			t.Skip("Skipping test - API might be unavailable")
		}

		// Check that we got valid data
		if poolInfo == nil {
			t.Error("Expected pool info, got nil")
		}

		// Check basic fields
		if poolInfo.ID != poolID {
			t.Errorf("Expected pool ID %s, got %s", poolID, poolInfo.ID)
		}

		// Check that we have token information
		if poolInfo.MintA.Address == "" || poolInfo.MintB.Address == "" {
			t.Error("Expected token addresses, got empty")
		}

		// Check that we have reserves
		if poolInfo.MintAmountA == 0 || poolInfo.MintAmountB == 0 {
			t.Error("Expected non-zero reserves")
		}

		// Log the pool info
		t.Logf("Pool: %s", poolInfo.ID)
		t.Logf("Token A: %s (%s) - Amount: %.2f", poolInfo.MintA.Symbol, poolInfo.MintA.Address, poolInfo.MintAmountA)
		t.Logf("Token B: %s (%s) - Amount: %.2f", poolInfo.MintB.Symbol, poolInfo.MintB.Address, poolInfo.MintAmountB)
		t.Logf("Price: %.6f", poolInfo.Price)
		t.Logf("TVL: $%.2f", poolInfo.Tvl)
		t.Logf("Fee Rate: %.4f%%", poolInfo.FeeRate*100)
	})

	// Test with invalid pool ID
	t.Run("Invalid pool ID", func(t *testing.T) {
		poolInfo, err := client.GetPoolInfoV3("InvalidPoolID123")
		if err == nil {
			t.Error("Expected error for invalid pool ID, got nil")
		}
		if poolInfo != nil {
			t.Error("Expected nil pool info for invalid ID")
		}
	})
}

func TestFindPoolByTokenV3(t *testing.T) {
	// Check if API is accessible
	resp, err := http.Get("https://api-v3.raydium.io/pools/info/mint?mint1=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v&poolType=all&poolSortField=liquidity&sortType=desc&pageSize=20&page=1")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Skip("Raydium API v3 is not accessible")
	}
	resp.Body.Close()

	client := NewAPIClient()

	tests := []struct {
		name      string
		tokenMint string
		wantErr   bool
	}{
		{
			name:      "Find USDC pool",
			tokenMint: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
			wantErr:   false,
		},
		{
			name:      "Find USDT pool",
			tokenMint: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", // USDT
			wantErr:   false,
		},
		{
			name:      "Unknown token",
			tokenMint: "UnknownTokenAddress123456789",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := client.FindPoolByTokenV3(tt.tokenMint)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindPoolByTokenV3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && pool != nil {
				// Verify pool has required fields
				if pool.ID == "" {
					t.Error("Pool ID is empty")
				}

				// Check that it's paired with SOL
				solMint := "So11111111111111111111111111111111111111112"
				isPairedWithSOL := strings.EqualFold(pool.MintA.Address, solMint) ||
					strings.EqualFold(pool.MintB.Address, solMint)

				if !isPairedWithSOL {
					t.Error("Pool is not paired with SOL")
				}

				// Verify pool has reserves
				if pool.MintAmountA <= 0 || pool.MintAmountB <= 0 {
					t.Error("Pool has zero reserves")
				}

				t.Logf("Found pool %s: %s-%s, TVL: $%.2f",
					pool.ID, pool.MintA.Symbol, pool.MintB.Symbol, pool.Tvl)
			}
		})
	}
}

func TestV3APIIntegration(t *testing.T) {
	client := NewAPIClient()

	t.Run("Find pool and get detailed info", func(t *testing.T) {
		// First find a pool for USDC
		pool, err := client.FindPoolByTokenV3("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
		if err != nil {
			t.Skip("API unavailable")
		}

		// Then get detailed info for the same pool
		detailedPool, err := client.GetPoolInfoV3(pool.ID)
		if err != nil {
			t.Fatalf("Failed to get detailed pool info: %v", err)
		}

		// Verify they match
		if detailedPool.ID != pool.ID {
			t.Errorf("Pool IDs don't match: %s != %s", detailedPool.ID, pool.ID)
		}

		// Both should have the same token pair
		if detailedPool.MintA.Address != pool.MintA.Address ||
			detailedPool.MintB.Address != pool.MintB.Address {
			t.Error("Token addresses don't match between search and detailed info")
		}

		t.Logf("Integration test passed: %s pool data is consistent", pool.ID)
	})
}
