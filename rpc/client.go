package rpc

import (
	"crypto/ecdsa"
	"net/http"
)

// RequestBody represents the JSON-RPC request body for mev_sendBundle API.
type RequestBody struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Client struct {
	httpClient *http.Client
	privKey    *ecdsa.PrivateKey
	baseURL    string
}

// NewClient creates a new instance of the API client.
func NewClient(clientURL string, auth ecdsa.PrivateKey) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    clientURL,
		privKey:    &auth,
	}
}
