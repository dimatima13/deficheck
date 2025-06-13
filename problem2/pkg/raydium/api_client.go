package raydium

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	apiV3BaseURL = "https://api-v3.raydium.io"
)

type APIClient struct {
	httpClient *http.Client
}

func NewAPIClient() *APIClient {
	return &APIClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type PoolInfoResponse struct {
	ID      string          `json:"id"`
	Success bool            `json:"success"`
	Data    []*PoolInfoData `json:"data"` // Use pointer to handle null values
}

type PoolInfoData struct {
	Type               string        `json:"type"`
	ProgramID          string        `json:"programId"`
	ID                 string        `json:"id"`
	MintA              TokenInfo     `json:"mintA"`
	MintB              TokenInfo     `json:"mintB"`
	Price              float64       `json:"price"`
	MintAmountA        float64       `json:"mintAmountA"`
	MintAmountB        float64       `json:"mintAmountB"`
	FeeRate            float64       `json:"feeRate"`
	OpenTime           string        `json:"openTime"`
	Tvl                float64       `json:"tvl"`
	Day                MarketData    `json:"day"`
	Week               MarketData    `json:"week"`
	Month              MarketData    `json:"month"`
	PoolType           []string      `json:"pooltype"`
	RewardDefaultInfos []interface{} `json:"rewardDefaultInfos"`
	FarmUpcomingCount  int           `json:"farmUpcomingCount"`
	FarmOngoingCount   int           `json:"farmOngoingCount"`
	FarmFinishedCount  int           `json:"farmFinishedCount"`
}

type TokenInfo struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	Address  string `json:"address"`
}

type MarketData struct {
	Volume      float64       `json:"volume"`
	VolumeQuote float64       `json:"volumeQuote"`
	VolumeFee   float64       `json:"volumeFee"`
	Apr         float64       `json:"apr"`
	FeeApr      float64       `json:"feeApr"`
	PriceMin    float64       `json:"priceMin"`
	PriceMax    float64       `json:"priceMax"`
	RewardApr   []interface{} `json:"rewardApr"`
}

func (c *APIClient) GetPoolInfoV3(poolID string) (*PoolInfoData, error) {
	url := fmt.Sprintf("%s/pools/info/ids?ids=%s", apiV3BaseURL, poolID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pool info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var poolResp PoolInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&poolResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !poolResp.Success || len(poolResp.Data) == 0 {
		return nil, fmt.Errorf("no pool data found for ID: %s", poolID)
	}

	// Check if the data is null (API returns [null] for invalid IDs)
	if poolResp.Data[0] == nil {
		return nil, fmt.Errorf("pool not found for ID: %s", poolID)
	}

	return poolResp.Data[0], nil
}

type PoolSearchResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Count       int              `json:"count"`
		Data        []PoolSearchData `json:"data"`
		HasNextPage bool             `json:"hasNextPage"`
	} `json:"data"`
}

type PoolSearchData struct {
	Type        string    `json:"type"`
	ProgramID   string    `json:"programId"`
	ID          string    `json:"id"`
	MintA       TokenInfo `json:"mintA"`
	MintB       TokenInfo `json:"mintB"`
	Price       float64   `json:"price"`
	Liquidity   float64   `json:"liquidity"`
	FeeRate     float64   `json:"feeRate"`
	MintAmountA float64   `json:"mintAmountA"`
	MintAmountB float64   `json:"mintAmountB"`
	TVL         float64   `json:"tvl"`
}

func (c *APIClient) FindPoolByTokenV3(tokenMint string) (*PoolInfoData, error) {
	solMint := "So11111111111111111111111111111111111111112"

	url := fmt.Sprintf("%s/pools/info/mint?mint1=%s&poolType=all&poolSortField=liquidity&sortType=desc&pageSize=20&page=1",
		apiV3BaseURL, tokenMint)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pools: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp PoolSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !searchResp.Success || searchResp.Data.Count == 0 {
		return nil, fmt.Errorf("no pools found for token %s", tokenMint)
	}

	var bestPoolID string
	var maxLiquidity float64

	for _, pool := range searchResp.Data.Data {
		isPairedWithSOL := strings.EqualFold(pool.MintA.Address, solMint) || strings.EqualFold(pool.MintB.Address, solMint)

		// Use TVL if Liquidity is 0 (some pools report TVL instead)
		poolLiquidity := pool.Liquidity
		if poolLiquidity == 0 {
			poolLiquidity = pool.TVL
		}

		if isPairedWithSOL && poolLiquidity > maxLiquidity {
			bestPoolID = pool.ID
			maxLiquidity = poolLiquidity
		}
	}

	if bestPoolID == "" {
		return nil, fmt.Errorf("no pool found for token %s paired with SOL", tokenMint)
	}

	return c.GetPoolInfoV3(bestPoolID)
}
