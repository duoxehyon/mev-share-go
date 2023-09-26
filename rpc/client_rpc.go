package rpc

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/metachris/flashbotsrpc"
)

// RequestBody represents the JSON-RPC request body for mev_sendBundle API.
type RequestBody struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// RPC client
type Client struct {
	httpClient      *http.Client
	flashbotsClient *flashbotsrpc.FlashbotsRPC
	privKey         *ecdsa.PrivateKey
	baseURL         string
}

// If request returns an error
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (err Error) Error() error {
	return fmt.Errorf("server returned error, Error Code: %d, Message: %s", err.Code, err.Message)
}

// Regular response
type Response[T any] struct {
	ID      int64  `json:"id"`
	Result  T      `json:"result"`
	Error   *Error `json:"error,omitempty"`
	JSONRPC string `json:"jsonrpc"`
}

// NewClient creates a new instance of the API client
func NewClient(clientURL string, auth *ecdsa.PrivateKey) MevAPIClient {
	flashbotsClient := flashbotsrpc.New(clientURL)

	return &Client{
		httpClient:      &http.Client{},
		flashbotsClient: flashbotsClient,
		baseURL:         clientURL,
		privKey:         auth,
	}
}

// Does api requests with Flashbots signature header
// returns the body
func (c *Client) CallWithSig(method string, params ...interface{}) ([]byte, error) {
	res, err := c.flashbotsClient.CallWithFlashbotsSignature(method, c.privKey, params...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Send private transaction ~`eth_sendPrivateTransaction`
// signedRawTx - transaction with nonce and vrs values
// options - options for private tx hints, builders, inclution, etc...
// returns the Transaction hash of the sent transaction
func (c *Client) SendPrivateTransaction(signedRawTx string, options *PrivateTxOptions) (*common.Hash, error) {
	tx := encodePrivateTxParams(signedRawTx, options)

	res, err := c.CallWithSig("eth_sendPrivateTransaction", tx)
	if err != nil {
		return nil, err
	}

	var decoded Response[common.Hash]
	err = json.Unmarshal(res, &decoded)
	if err != nil {
		return nil, err
	}

	if decoded.Error != nil {
		return nil, decoded.Error.Error()
	}

	return &decoded.Result, nil
}

// Send mev-share bundle  ~`mev_sendBundle`
// bundle - the bundle with all transactions / hashes
// returns the bundle hash / error
func (c *Client) SendBundle(bundle MevSendBundleParams) (*MevSendBundleResponse, error) {
	bundle.Version = "v0.1"
	res, err := c.CallWithSig("mev_sendBundle", bundle)
	if err != nil {
		return nil, err
	}

	var decoded Response[MevSendBundleResponse]
	err = json.Unmarshal(res, &decoded)
	if err != nil {
		return nil, err
	}

	if decoded.Error != nil {
		return nil, decoded.Error.Error()
	}

	return &decoded.Result, nil
}

// Simulate bundle ~`mev_simBundle`
// bundle - the bundle with all transactions / hashes
// simOverrides - given values will be overwritten when doing the simulation
// returns the simulation result / error
func (c *Client) SimBundle(bundle MevSendBundleParams, simOverrides SimBundleOverrides) (*SimBundleResponse, error) {
	bundle.Version = "v0.1"
	res, err := c.CallWithSig("mev_simBundle", bundle, simOverrides)
	if err != nil {
		return nil, err
	}

	var decoded Response[SimBundleResponse]
	err = json.Unmarshal(res, &decoded)
	if err != nil {
		return nil, err
	}

	if decoded.Error != nil {
		return nil, decoded.Error.Error()
	}

	return &decoded.Result, nil
}
