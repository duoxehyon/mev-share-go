package shared

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

// Custom type because of hex string to bytes decoding error while using default geth.Log
type Log struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data,omitempty"` // Could be replaces with geth.hexutil type
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
