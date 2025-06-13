# Problem 2: Raydium DeFi Quote Service

A Solana-based quote service that interacts with Raydium AMM protocol to provide real-time price quotes for token swaps.

## Overview

This solution implements a quote service for Raydium, one of the leading AMM protocols on Solana. The program retrieves pool data directly from the blockchain and calculates swap prices using the constant product formula (x * y = k).

## How to Run

### Build
```bash
go build -o problem2 cmd/main.go
```

### Basic Usage
```bash
# Get a quote for buying USDC with 100 SOL
./problem2 -token EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v -qty 100 -side buy

# Get a quote for selling USDT for SOL
./problem2 -token Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB -qty 50 -side sell
```

### Advanced Options
- `-api`: Enable dynamic pool discovery via Raydium API
- `-onchain`: Fetch all data directly from blockchain (slower but real-time)
- `-mock`: Use mock data for testing
- `-rpc <URL>`: Use custom Solana RPC endpoint

## Architecture & Design Decisions

### Data Fetching Strategy
1. **Hardcoded pools** (default): Fastest option for known tokens (USDC, USDT)
2. **API mode**: Uses Raydium v3 API for dynamic pool discovery
3. **Onchain mode**: Direct blockchain queries for real-time accuracy

I chose to optimize for **time complexity** over space because:
- DeFi applications require low latency for accurate quotes
- Pool data is relatively small and caching improves response times
- Users expect real-time prices in trading applications

### Performance Complexity

**Time Complexity:**
- Hardcoded mode: O(1) - Direct pool lookup
- API mode: O(n) where n is pools returned by search (~20-100)
- Onchain mode: O(1) for pool data + RPC latency

**Space Complexity:**
- O(p) where p is number of cached pools
- Minimal memory footprint (~100KB per pool)

**Optimizations Implemented:**
- Pool caching to avoid repeated API/RPC calls
- Direct v3 API integration for efficient pool search
- Batch RPC requests where possible

### Key Design Choices

1. **Direct Smart Contract Interaction**: Instead of relying solely on APIs, the solution can parse Raydium V4 pool data directly from the blockchain, ensuring accuracy and decentralization.

2. **Custom Base58 Implementation**: Implemented our own Base58 encoding to handle Solana addresses without external dependencies.

3. **Big Number Arithmetic**: Used Go's `math/big` for precise calculations, critical for DeFi applications.

## External Libraries

- **Standard Go libraries only**: No external dependencies beyond Go's standard library
- This decision was made to:
  - Minimize security risks
  - Ensure easy deployment
  - Maintain full control over the codebase

## Protocol Choice: Raydium

Raydium was selected because:
- Open-source AMM with well-documented pool structure
- High liquidity on Solana
- Standard constant product formula (x * y = k)
- Direct blockchain data access without proprietary APIs

## AI Assistance Disclosure

During development, I used Claude (Anthropic) as a coding assistant for:
- **Code review and optimization suggestions**: Helped identify performance improvements like transitioning from v2 to v3 API
- **Documentation assistance**: Helped structure comments and README sections
- **Testing strategies**: Suggested edge cases for unit tests
- **Debugging support**: Assisted in troubleshooting API response parsing

**Validation steps taken:**
1. Manually tested all code paths with real Solana mainnet data
2. Verified AMM calculations against Raydium's actual swap amounts
3. Cross-referenced pool data with Solana explorers
4. Implemented comprehensive unit tests for critical functions
5. Verified v3 API functionality with comprehensive testing

The core architecture, AMM formula implementation, and smart contract parsing logic were designed based on my understanding of DeFi protocols and Solana's architecture. AI was primarily used as a productivity tool for faster iteration and catching potential issues.

## Testing

Run tests:
```bash
go test ./...
```

Run specific test:
```bash
go test ./pkg/raydium -v
```

## Future Improvements

1. Support for concentrated liquidity pools
2. Multi-hop routing for better prices
3. WebSocket support for real-time price updates
4. Integration with more Solana DeFi protocols