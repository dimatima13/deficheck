package solana

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	rpcURL string
	client *http.Client
}

func NewClient(rpcURL string) *Client {
	if rpcURL == "" {
		rpcURL = "https://api.mainnet-beta.solana.com"
	}
	return &Client{
		rpcURL: rpcURL,
		client: &http.Client{},
	}
}

type RPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      int             `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) GetAccountInfo(address string) (map[string]interface{}, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  "getAccountInfo",
		Params: []interface{}{
			address,
			map[string]string{"encoding": "base64"},
		},
		ID: 1,
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

func (c *Client) GetMultipleAccounts(addresses []string) ([]map[string]interface{}, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  "getMultipleAccounts",
		Params: []interface{}{
			addresses,
			map[string]string{"encoding": "base64"},
		},
		ID: 1,
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result.Value, nil
}

func (c *Client) GetTokenAccountBalance(address string) (map[string]interface{}, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  "getTokenAccountBalance",
		Params:  []interface{}{address},
		ID:      1,
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

func (c *Client) doRequest(req RPCRequest) (*RPCResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.rpcURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s (code: %d)", rpcResp.Error.Message, rpcResp.Error.Code)
	}

	return &rpcResp, nil
}

func DecodeBase64Data(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
