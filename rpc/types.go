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
