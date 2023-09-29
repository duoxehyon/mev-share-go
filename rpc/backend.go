// Package rpc is the RPC Client for `mev` api
package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/flashbots/mev-share-node/mevshare"
)

// MevAPIClient is the MEV-Share Client abstraction
type MevAPIClient interface {
	// MEV-Share Api Requests with Flashbots signature header
	CallWithSig(method string, params ...interface{}) ([]byte, error)
	// Send bundle tom mev-share node
	SendBundle(bundle mevshare.SendMevBundleArgs) (*mevshare.SendMevBundleResponse, error)
	// Bundle simulation
	SimBundle(bundle mevshare.SendMevBundleArgs, simOverrides mevshare.SimMevBundleAuxArgs) (*mevshare.SimMevBundleResponse, error)
	// Send private transaction with hints
	SendPrivateTransaction(signedRawTx string, options *PrivateTxOptions) (*common.Hash, error)
}
