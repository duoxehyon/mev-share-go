package sse

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Event represents a matchmaker event sent from sse subscription.
type Event struct {
	Data  *MatchMakerEvent // Will be nil if an error occurred during poll
	Error error
}

// MatchMakerEvent represents the pending transaction hints sent by matchmaker.
type MatchMakerEvent struct {
	Hash        common.Hash          `json:"hash"`
	Logs        []Log                `json:"logs,omitempty"`
	Txs         []PendingTransaction `json:"txs,omitempty"`
	MevGasPrice *big.Int             `json:"mevGasPrice,omitempty"`
	GasUsed     *big.Int             `json:"gasUsed,omitempty"`
}

// PendingTransaction represents the hits revealed by the matchmaker about the transaction / bundle.
type PendingTransaction struct {
	To               common.Address `json:"to"`
	FunctionSelector [4]byte        `json:"functionSelector,omitempty"`
	CallData         []byte         `json:"callData,omitempty"`
	MevGasPrice      *big.Int       `json:"mevGasPrice,omitempty"`
	GasUsed          *big.Int       `json:"gasUsed,omitempty"`
}

// UnmarshalJSON unmarshals JSON data into a PendingTransaction.
func (t *PendingTransaction) UnmarshalJSON(data []byte) error {
	var temp struct {
		To               common.Address `json:"to"`
		FunctionSelector string         `json:"functionSelector,omitempty"`
		CallData         string         `json:"callData,omitempty"`
		MevGasPrice      *big.Int       `json:"mevGasPrice,omitempty"`
		GasUsed          *big.Int       `json:"gasUsed,omitempty"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	t.To = temp.To
	t.MevGasPrice = temp.MevGasPrice
	t.GasUsed = temp.GasUsed

	if temp.CallData != "" && temp.CallData != "0x" {
		decoded, err := hex.DecodeString(strings.TrimPrefix(temp.CallData, "0x"))
		if err == nil {
			t.CallData = decoded
		}
	}

	if temp.FunctionSelector != "" && temp.FunctionSelector != "0x" {
		decoded, err := hex.DecodeString(strings.TrimPrefix(temp.FunctionSelector, "0x"))
		if err == nil && len(decoded) >= 4 {
			copy(t.FunctionSelector[:], decoded[:4])
		}
	}

	return nil
}

// Custom type because of hex string to bytes decoding error while using default geth.Log
type Log struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data,omitempty"`
}

// UnmarshalJSON unmarshals JSON data into a Log.
func (l *Log) UnmarshalJSON(data []byte) error {
	var temp struct {
		Address common.Address `json:"address"`
		Topics  []common.Hash  `json:"topics"`
		Data    string         `json:"data,omitempty"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	l.Address = temp.Address
	l.Topics = temp.Topics

	if temp.Data != "" {
		if strings.HasPrefix(temp.Data, "0x") {
			temp.Data = temp.Data[2:]
		}

		decoded, err := hex.DecodeString(temp.Data)
		if err != nil {
			return err
		}

		l.Data = decoded
	} else {
		l.Data = nil
	}

	return nil
}
