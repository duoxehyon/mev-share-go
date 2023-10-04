package rpc

import "github.com/ethereum/go-ethereum/common/hexutil"

// Regular transaction
type SignedRawTx struct {
	Tx        string `json:"tx,omitempty"`
	CanRevert bool   `json:"canRevert,omitempty"`
}

// Hints represents hints for privacy preferences
type Hints struct {
	CallData         bool
	ContractAddress  bool
	FunctionSelector bool
	Logs             bool
	DefaultLogs      bool
	TxHash           bool
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
	if h.DefaultLogs {
		hints = append(hints, "default_logs")
	}
	if h.TxHash {
		hints = append(hints, "tx_hash")
	}
	// hints = append(hints, "hash")

	return hints
}

// `eth_sendPrivateTransaction` parameters
type PrivateTxOptions struct {
	Hints          Hints
	MaxBlockNumber hexutil.Uint64
	Builders       []string
}

// Encodes data for mev_sendPrivateTransaction
func encodePrivateTxParams(signedTx string, options *PrivateTxOptions) interface{} {
	// Create the Privacy struct
	privacy := struct {
		Hints []string `json:"hints,omitempty"`
	}{
		Hints: options.Hints.String(),
	}

	preferences := struct {
		Fast     bool        `json:"fast"`
		Privacy  interface{} `json:"privacy,omitempty"`
		Builders []string    `json:"builders,omitempty"`
	}{
		Fast:     true,
		Privacy:  privacy,
		Builders: options.Builders,
	}

	data := struct {
		Tx             string         `json:"tx"`
		MaxBlockNumber hexutil.Uint64 `json:"maxBlockNumber,omitempty"`
		Preferences    interface{}    `json:"preferences"`
	}{
		Tx:             signedTx,
		MaxBlockNumber: options.MaxBlockNumber,
		Preferences:    preferences,
	}

	return data
}
