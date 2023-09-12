package shared

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
)

// Log - Custom type because of hex string to bytes decoding error while using default geth.Log
type Log struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data,omitempty"` // Could be replaced with geth.hexutil type
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

	l.Topics = temp.Topics
	l.Address = temp.Address

	if temp.Data != "" {
		temp.Data = temp.Data[2:]

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
