package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"deficheck/problem2/internal/types"
	"deficheck/problem2/internal/quote"
	"deficheck/problem2/pkg/raydium"
	"deficheck/problem2/pkg/solana"
)

func main() {
	var (
		tokenAddress = flag.String("token", "", "Token contract address")
		quantity     = flag.String("qty", "", "Quantity to trade")
		side         = flag.String("side", "", "Trade side: buy or sell")
		rpcURL       = flag.String("rpc", "https://api.mainnet-beta.solana.com", "Solana RPC URL")
		mockMode     = flag.Bool("mock", false, "Use mock data instead of real blockchain data")
		useAPI       = flag.Bool("api", false, "Use Raydium API to find pools dynamically")
		useOnchain   = flag.Bool("onchain", false, "Fetch all data directly from blockchain (fully onchain)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -token <address> -qty <amount> -side <buy|sell>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Get a price quote from Raydium DEX on Solana\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -token EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v -qty 100 -side buy\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNote: Uses Raydium protocol. All pairs denominated in SOL/wSOL.\n")
	}

	flag.Parse()

	// Validate required flags
	if *tokenAddress == "" || *quantity == "" || *side == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Parse quantity
	qty, ok := new(big.Float).SetString(*quantity)
	if !ok || qty.Sign() <= 0 {
		log.Fatalf("Invalid quantity: %s", *quantity)
	}

	// Validate side
	*side = strings.ToLower(*side)
	if *side != "buy" && *side != "sell" {
		log.Fatalf("Invalid side: %s (must be 'buy' or 'sell')", *side)
	}

	// Create clients
	solanaClient := solana.NewClient(*rpcURL)
	raydiumClient := raydium.NewClient(solanaClient)
	quoteService := quote.NewService(raydiumClient)
	
	// Enable API mode if requested
	if *useAPI {
		quoteService.SetUseAPI(true)
		fmt.Println("API mode enabled - will search for pools dynamically")
	}
	
	// Enable onchain mode if requested
	if *useOnchain {
		quoteService.SetUseOnchain(true)
		fmt.Println("Onchain mode enabled - will fetch all data from blockchain")
	}

	// Create quote request
	request := &types.QuoteRequest{
		TokenAddress: *tokenAddress,
		Quantity:     qty,
		Side:         *side,
	}

	// Get quote
	fmt.Printf("Fetching quote from %s...\n", types.ProtocolName)
	
	var quoteResult *types.QuoteResponse
	var err error
	
	if *mockMode {
		fmt.Println("Using mock data...")
		// Use mock pool for testing
		mockPool, err := quote.GetMockPool(request.TokenAddress)
		if err != nil {
			log.Fatalf("Failed to get mock pool: %v", err)
		}
		
		price, err := raydiumClient.CalculatePrice(mockPool, request.Quantity, request.Side)
		if err != nil {
			log.Fatalf("Failed to calculate price: %v", err)
		}
		
		decimals := quote.DetermineDecimals(price)
		priceFormatted := quote.FormatPrice(price, decimals)
		
		quoteResult = &types.QuoteResponse{
			Price:          price,
			PriceFormatted: priceFormatted,
			TokenSymbol:    quote.GetTokenSymbol(request.TokenAddress),
			Decimals:       decimals,
			Protocol:       types.ProtocolName,
		}
	} else {
		quoteResult, err = quoteService.GetQuote(request)
		if err != nil {
			log.Fatalf("Failed to get quote: %v", err)
		}
	}

	// Display results
	fmt.Println("\n===== QUOTE RESULT =====")
	fmt.Printf("Protocol: %s\n", quoteResult.Protocol)
	fmt.Printf("Token: %s\n", quoteResult.TokenSymbol)
	fmt.Printf("Side: %s\n", *side)
	fmt.Printf("Quantity: %s\n", *quantity)
	fmt.Printf("Price: %s SOL\n", quoteResult.PriceFormatted)
	fmt.Printf("Decimals: %d\n", quoteResult.Decimals)
	fmt.Println("=======================")
}
