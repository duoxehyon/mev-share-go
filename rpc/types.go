package rpc

import (
	"encoding/json"
	"github.com/duoxehyon/mev-share-go/shared"

	"github.com/ethereum/go-ethereum/common"
)

// MevSendBundleParams represents the parameters for mev_sendBundle
type MevSendBundleParams struct {
	Version   string        `json:"version"`
	Inclusion Inclusion     `json:"inclusion"`
	Body      []BundleItem  `json:"body"`
	Validity  Validity      `json:"validity,omitempty"`
	Privacy   BundlePrivacy `json:"privacy,omitempty"`
}

// Inclusion represents the inclusion data for mev_sendBundle
type Inclusion struct {
	Block    uint64  `json:"block"`
	MaxBlock *uint64 `json:"maxBlock,omitempty"`
}

// Custom interface that both SignedRawTx and string types satisfy
type BundleItem interface {
	isTxType()
}

func (t MevShareTxHash) isTxType() {}
func (t SignedRawTx) isTxType()    {}

type MevShareTxHash struct {
	Hash string
}

type SignedRawTx struct {
	Tx        string `json:"tx,omitempty"`
	CanRevert bool   `json:"canRevert,omitempty"`
}

// Validity represents the validity conditions for the bundle.
type Validity struct {
	Refund       []Refund       `json:"refund,omitempty"`
	RefundConfig []RefundConfig `json:"refundConfig,omitempty"`
}

// RefundConfig represents the refund configuration for the bundle.
type RefundConfig struct {
	BodyIdx int `json:"bodyIdx"`
	Percent int `json:"percent"`
}

type Refund struct {
	Address common.Address `json:"address"`
	Percent int            `json:"percent"`
}

// Hints represents hints for privacy preferences.
type Hints struct {
	CallData         bool
	ContractAddress  bool
	FunctionSelector bool
	Logs             bool
	TxHash           bool
}

func (h *Hints) String() []string {
	hints := make([]string, 0)
	if h.CallData {
		hints = append(hints, "calldata")
	}
	if h.ContractAddress {
		hints = append(hints, "contractAddress")
	}
	if h.FunctionSelector {
		hints = append(hints, "functionSelector")
	}
	if h.Logs {
		hints = append(hints, "logs")
	}
	if h.TxHash {
		hints = append(hints, "txHash")
	}

	return hints
}

// BundlePrivacy represents the privacy preferences for the bundle.
type BundlePrivacy struct {
	Hints    Hints    `json:"-"`
	Builders []string `json:"builders,omitempty"`
}

// MarshalJSON implements the custom JSON marshaling for BundlePrivacy.
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

type PrivateTxOptions struct {
	Hints          Hints
	MaxBlockNumber uint64
	Builders       []string
}

// Encodes data for mev_sendPrivateTransaction
func encodePrivateTxParams(signedTx string, options *PrivateTxOptions) ([]interface{}, error) {
	data := struct {
		Tx             string
		MaxBlockNumber uint64
		Preferences    struct {
			Fast     bool
			Privacy  []string
			Builders []string
		}
	}{
		Tx:             signedTx,
		MaxBlockNumber: options.MaxBlockNumber,
		Preferences: struct {
			Fast     bool
			Privacy  []string
			Builders []string
		}{
			Fast:     true,
			Privacy:  options.Hints.String(),
			Builders: options.Builders,
		},
	}

	return []interface{}{data}, nil
}

// Bundle simulation parameters for mev_simBundle
type SimBundleOverrides struct {
	// Block used for simulation state. Defaults to latest block
	ParentBlock uint64 `json:"parentBlock,omitempty"`
	// Block number used for simulation, defaults to parentBlock.number + 1
	BlockNumber uint64 `json:"blockNumber,omitempty"`
	// Coinbase used for simulation, defaults to parentBlock.coinbase
	Coinbase uint64 `json:"coinbase,omitempty"`
	// Timestamp used for simulation, defaults to parentBlock.timestamp + 12
	Timestamp uint64 `json:"timestamp,omitempty"`
	// Gas limit used for simulation, defaults to parentBlock.gasLimit
	GasLimit uint64 `json:"gasLimit,omitempty"`
	// Base fee used for simulation, defaults to parentBlock.baseFeePerGas
	BaseFee uint64 `json:"baseFee,omitempty"`
	// Timeout in seconds, defaults to 5
	Timeout uint64 `json:"timeout,omitempty"`
}

// Response for mev_simBundle
type SimBundleResponse struct {
	/// Whether the simulation was successful
	Success bool `json:"success"`
	// Error if simulation failed
	Error *string `json:"error,omitempty"`
	// The block number of the simulated block
	StateBlock uint64 `json:"stateBlock"`
	// The gas price of the simulated block
	MevGasPrice uint64 `json:"mevGasPrice"`
	// The profit of the simulated block
	Profit uint64 `json:"profit"`
	// The refundable value of the simulated block
	RefundableValue uint64 `json:"refundableValue"`
	// The gas used by the simulated block.
	GasUsed uint64 `json:"gasUsed"`
	// Logs returned by mev_simBundle.
	Logs *[]SimBundleLogs `json:"logs,omitempty"`
}

// Logs returned by mev_simBundle
type SimBundleLogs struct {
	// Logs from transactions
	TxLogs *[]shared.Log `json:"txLogs,omitempty"`
	// Logs for bundles inside bundle
	BundleLogs *[]SimBundleLogs `json:"bundleLogs,omitempty"`
}
