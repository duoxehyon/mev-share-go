package rpc

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
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

func (c *Client) SendPrivateTransaction(signedRawTx string, options PrivateTxOptions) (string, error) {
	tx, err := encodePrivateTxParams(signedRawTx, &options)
	if err != nil {
		return "", err
	}

	res, err := c.CallWithSig("eth_sendPrivateTransaction", tx)
	if err != nil {
		return "", err
	}

	var hash string
	err = json.Unmarshal(res, &hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}
