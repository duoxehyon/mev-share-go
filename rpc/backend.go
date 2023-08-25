package rpc

import "github.com/ethereum/go-ethereum/common"

// MEV-Share Client abstraction
type MevApiClient interface {
	// MEV-Share Api Requests with Flashbots signature header
	CallWithSig(method string, params ...interface{}) ([]byte, error)
	// Send bundle tom mev-share node
	SendBundle(bundle MevSendBundleParams) (*MevSendBundleResponse, error)
	// Bundle simulation
	SimBundle(bundle MevSendBundleParams, simOverrides SimBundleOverrides) (*SimBundleResponse, error)
	// Send private transaction with hints
	SendPrivateTransaction(signedRawTx string, options *PrivateTxOptions) (*common.Hash, error)
}
