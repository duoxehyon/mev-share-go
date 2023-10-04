package rpc

import (
	"crypto/ecdsa"
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/flashbots/mev-share-node/mevshare"
	"github.com/metachris/flashbotsrpc"
)

// RPC client
type Client struct {
	httpClient      *http.Client
	flashbotsClient *flashbotsrpc.FlashbotsRPC
	privKey         *ecdsa.PrivateKey
	baseURL         string
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
	return c.flashbotsClient.CallWithFlashbotsSignature(method, c.privKey, params...)
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

	var decoded common.Hash
	err = json.Unmarshal(res, &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}

type (
	SendMevBundleArgs     = mevshare.SendMevBundleArgs
	SendMevBundleResponse = mevshare.SendMevBundleResponse
	MevBundleBody         = mevshare.MevBundleBody
	Inclusion             = mevshare.MevBundleInclusion
	Validity              = mevshare.MevBundleValidity
	Refund                = mevshare.RefundConstraint
	RefundConfig          = mevshare.RefundConfig
	Privacy               = mevshare.MevBundlePrivacy
	HintIntent            = mevshare.HintIntent
	MetaData              = mevshare.MevBundleMetadata
)

// Send mev-share bundle  ~`mev_sendBundle`
// bundle - the bundle with all transactions / hashes
// returns the bundle hash / error
func (c *Client) SendBundle(bundle SendMevBundleArgs) (*mevshare.SendMevBundleResponse, error) {
	bundle.Version = "v0.1"
	res, err := c.CallWithSig("mev_sendBundle", bundle)
	if err != nil {
		return nil, err
	}

	var decoded mevshare.SendMevBundleResponse
	err = json.Unmarshal(res, &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}

type (
	SimMevBundleAuxArgs  = mevshare.SimMevBundleAuxArgs
	SimMevBundleResponse = mevshare.SimMevBundleResponse
)

// Simulate bundle ~`mev_simBundle`
// bundle - the bundle with all transactions / hashes
// simOverrides - given values will be overwritten when doing the simulation
// returns the simulation result / error
func (c *Client) SimBundle(bundle mevshare.SendMevBundleArgs, simOverrides mevshare.SimMevBundleAuxArgs) (*mevshare.SimMevBundleResponse, error) {
	bundle.Version = "v0.1"
	res, err := c.CallWithSig("mev_simBundle", bundle, simOverrides)
	if err != nil {
		return nil, err
	}

	var decoded mevshare.SimMevBundleResponse
	err = json.Unmarshal(res, &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}
