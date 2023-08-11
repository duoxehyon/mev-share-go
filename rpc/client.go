package rpc

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
	httpClient *http.Client
	privKey    *ecdsa.PrivateKey
	baseURL    string
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (err Error) Error() error {
	return errors.New(fmt.Sprintf("Server Returned error, Error Code: %d, Message: %s", err.Code, err.Message))
}

type Response[T any] struct {
	ID      int64  `json:"id"`
	Result  T      `json:"result"`
	Error   *Error `json:"error,omitempty"`
	JSONRPC string `json:"jsonrpc"`
}

// NewClient creates a new instance of the API client.
func NewClient(clientURL string, auth *ecdsa.PrivateKey) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    clientURL,
		privKey:    auth,
	}
}

// Does api requests with Flashbots signature header
func (c *Client) CallWithSig(method string, params ...interface{}) ([]byte, error) {
	request := RequestBody{
		ID:      777,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	hashedBody := crypto.Keccak256Hash([]byte(body)).Hex()
	sig, err := crypto.Sign(accounts.TextHash([]byte(hashedBody)), c.privKey)
	if err != nil {
		return nil, err
	}

	signature := crypto.PubkeyToAddress(c.privKey.PublicKey).Hex() + ":" + hexutil.Encode(sig)

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Flashbots-Signature", signature)

	response, err := c.httpClient.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Send private transaction `eth_sendPrivateTransaction`
func (c *Client) SendPrivateTransaction(signedRawTx string, options PrivateTxOptions) (*common.Hash, error) {
	tx := encodePrivateTxParams(signedRawTx, &options)

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

// Send mev-share bundle `mev_sendBundle`
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

// Simulate bundle `mev_simBundle`
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
