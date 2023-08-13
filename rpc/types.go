package rpc

import (
	"encoding/json"

	"github.com/duoxehyon/mev-share-go/shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// MevSendBundleParams represents the parameters for mev_sendBundle
type MevSendBundleParams struct {
	Version   string         `json:"version"`
	Inclusion Inclusion      `json:"inclusion"`
	Body      []BundleItem   `json:"body"`
	Validity  *Validity      `json:"validity,omitempty"`
	Privacy   *BundlePrivacy `json:"privacy,omitempty"`
}

// Inclusion represents the inclusion data for mev_sendBundle
type Inclusion struct {
	Block    hexutil.Uint64 `json:"block"`
	MaxBlock hexutil.Uint64 `json:"maxBlock,omitempty"`
}

// Custom interface that both SignedRawTx and string types satisfy
type BundleItem interface {
	isTxType()
}

func (t MevShareTxHash) isTxType() {}
func (t SignedRawTx) isTxType()    {}

// Mev-share transaction hash
type MevShareTxHash struct {
	Hash string
}

// Regular transaction
type SignedRawTx struct {
	Tx        string `json:"tx,omitempty"`
	CanRevert bool   `json:"canRevert,omitempty"`
}

// Validity represents the validity conditions for the bundle
type Validity struct {
	Refund       []Refund       `json:"refund,omitempty"`
	RefundConfig []RefundConfig `json:"refundConfig,omitempty"`
}

// RefundConfig represents the refund configuration for the bundle
type RefundConfig struct {
	BodyIdx int `json:"bodyIdx"`
	Percent int `json:"percent"`
}

// Refund address and percentage
type Refund struct {
	Address common.Address `json:"address"`
	Percent int            `json:"percent"`
}

// Hints represents hints for privacy preferences
type Hints struct {
	CallData         bool
	ContractAddress  bool
	FunctionSelector bool
	Logs             bool
	TxHash           bool
	Hash             bool
	SpecialLogs      bool
}

func (h *Hints) String() []string {
	hints := make([]string, 0)
	if h.CallData {
		hints = append(hints, "calldata")
	}
	if h.ContractAddress {
		hints = append(hints, "contract_address")
	}
	if h.FunctionSelector {
		hints = append(hints, "function_selector")
	}
	if h.Logs {
		hints = append(hints, "logs")
	}
	if h.TxHash {
		hints = append(hints, "tx_hash")
	}
	if h.Hash {
		hints = append(hints, "hash")
	}
	if h.SpecialLogs {
		hints = append(hints, "special_logs")
	}

	return hints
}

// BundlePrivacy represents the privacy preferences for the bundle
type BundlePrivacy struct {
	Hints    Hints    `json:"-"`
	Builders []string `json:"builders,omitempty"`
}

// MarshalJSON implements the custom JSON marshaling for BundlePrivacy
func (bp BundlePrivacy) MarshalJSON() ([]byte, error) {
	hints := bp.Hints.String()

	data := struct {
		Hints    []string `json:"hints,omitempty"`
		Builders []string `json:"builders,omitempty"`
	}{
		Hints:    hints,
		Builders: bp.Builders,
	}

	return json.Marshal(data)
}

// Response for mev_sendBundle
type MevSendBundleResponse struct {
	BundleHash common.Hash `json:"bundleHash"`
}

// `eth_sendPrivateTransaction` parameters
type PrivateTxOptions struct {
	Hints          Hints
	MaxBlockNumber hexutil.Uint64
	Builders       []string
}

// Encodes data for mev_sendPrivateTransaction
func encodePrivateTxParams(signedTx string, options *PrivateTxOptions) any {
	data := struct {
		Tx             string         `json:"tx"`
		MaxBlockNumber hexutil.Uint64 `json:"maxBlockNumber,omitempty"`
		Preferences    struct {
			Fast     bool     `json:"fast"`
			Privacy  []string `json:"privacy,omitempty"`
			Builders []string `json:"builders,omitempty"`
		} `json:"preferences"`
	}{
		Tx:             signedTx,
		MaxBlockNumber: options.MaxBlockNumber,
		Preferences: struct {
			Fast     bool     `json:"fast"`
			Privacy  []string `json:"privacy,omitempty"`
			Builders []string `json:"builders,omitempty"`
		}{
			Fast:     true,
			Privacy:  options.Hints.String(),
			Builders: options.Builders,
		},
	}

	return data
}

// Bundle simulation parameters for mev_simBundle
type SimBundleOverrides struct {
	// Block used for simulation state. Defaults to latest block
	ParentBlock hexutil.Uint64 `json:"parentBlock,omitempty"`
	// Block number used for simulation, defaults to parentBlock.number + 1
	BlockNumber *hexutil.Big `json:"blockNumber,omitempty"`
	// Coinbase used for simulation, defaults to parentBlock.coinbase
	Coinbase *common.Address `json:"coinbase,omitempty"`
	// Timestamp used for simulation, defaults to parentBlock.timestamp + 12
	Timestamp hexutil.Uint64 `json:"timestamp,omitempty"`
	// Gas limit used for simulation, defaults to parentBlock.gasLimit
	GasLimit hexutil.Uint64 `json:"gasLimit,omitempty"`
	// Base fee used for simulation, defaults to parentBlock.baseFeePerGas
	BaseFee *hexutil.Big `json:"baseFee,omitempty"`
	// Timeout in seconds, defaults to 5
	Timeout uint64 `json:"timeout,omitempty"`
}

// Response for mev_simBundle
type SimBundleResponse struct {
	Success bool `json:"success"`
	// Error if simulation failed
	Error string `json:"error,omitempty"`
	// The block number of the simulated block
	StateBlock hexutil.Uint64 `json:"stateBlock"`
	// The gas price of the simulated block
	MevGasPrice hexutil.Big `json:"mevGasPrice"`
	// The profit of the simulated block
	Profit hexutil.Big `json:"profit"`
	// The refundable value of the simulated block
	RefundableValue hexutil.Big `json:"refundableValue"`
	// The gas used by the simulated block
	GasUsed hexutil.Uint64 `json:"gasUsed"`
	// Logs returned by mev_simBundle
	Logs *[]SimBundleLogs `json:"logs,omitempty"`
}

// Logs returned by mev_simBundle
type SimBundleLogs struct {
	// Logs inside transactions
	TxLogs *[]shared.Log `json:"txLogs,omitempty"`
	// Logs for bundles inside bundle
	BundleLogs *[]SimBundleLogs `json:"bundleLogs,omitempty"`
}
