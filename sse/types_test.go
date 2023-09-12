package sse

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestPendingTransaction_UnmarshalJSON(t *testing.T) {
	rawJSON := `{
		"to": "0x1234567890abcdef1234567890abcdef12345678",
		"functionSelector": "0xabcdef12",
		"callData": "0xdeadbeef",
		"mevGasPrice": "0x1000000000000",
		"gasUsed": "0x200000"
	}`

	expectedTo := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedFunctionSelector := [4]byte{0xab, 0xcd, 0xef, 0x12}
	expectedCallData := []byte{0xde, 0xad, 0xbe, 0xef}
	expectedMevGasPrice := big.NewInt(0x1000000000000)
	expectedGasUsed := big.NewInt(0x200000)

	var tx PendingTransaction
	err := json.Unmarshal([]byte(rawJSON), &tx)
	assert.NoError(t, err)
	assert.Equal(t, expectedTo, tx.To)
	assert.Equal(t, expectedFunctionSelector, tx.FunctionSelector)
	assert.Equal(t, expectedCallData, tx.CallData)
	assert.Equal(t, expectedMevGasPrice.String(), tx.MevGasPrice.ToInt().String())
	assert.Equal(t, expectedGasUsed.String(), tx.GasUsed.ToInt().String())
}

func TestPendingTransaction_UnmarshalJSON_EmptyFields(t *testing.T) {
	rawJSON := `{
		"to": "0x1234567890abcdef1234567890abcdef12345678",
		"functionSelector": "",
		"callData": "",
		"mevGasPrice": "0x0",
		"gasUsed": "0x0"
	}`

	expectedTo := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedFunctionSelector := [4]byte{}
	var expectedCallData []byte
	expectedMevGasPrice := big.NewInt(0)
	expectedGasUsed := big.NewInt(0)

	var tx PendingTransaction
	err := json.Unmarshal([]byte(rawJSON), &tx)
	assert.NoError(t, err)
	assert.Equal(t, expectedTo, tx.To)
	assert.Equal(t, expectedFunctionSelector, tx.FunctionSelector)
	assert.Equal(t, expectedCallData, tx.CallData)
	assert.Equal(t, expectedMevGasPrice.String(), tx.MevGasPrice.ToInt().String())
	assert.Equal(t, expectedGasUsed.String(), tx.GasUsed.ToInt().String())
}

func TestPendingTransaction_UnmarshalJSON_MissingFields(t *testing.T) {
	rawJSON := `{
		"to": "0x1234567890abcdef1234567890abcdef12345678"
	}`

	expectedTo := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	var expectedFunctionSelector [4]byte
	var expectedCallData []byte
	var expectedMevGasPrice *big.Int
	var expectedGasUsed *big.Int

	var tx PendingTransaction
	err := json.Unmarshal([]byte(rawJSON), &tx)
	assert.NoError(t, err)
	assert.Equal(t, expectedTo, tx.To)
	assert.Equal(t, expectedFunctionSelector, tx.FunctionSelector)
	assert.Equal(t, expectedCallData, tx.CallData)
	assert.Equal(t, expectedMevGasPrice.String(), tx.MevGasPrice.ToInt().String())
	assert.Equal(t, expectedGasUsed.String(), tx.GasUsed.ToInt().String())
}

func TestEvent_Data_Error(t *testing.T) {
	errMsg := "some error"
	event := Event{Error: errors.New(errMsg)}

	assert.Error(t, event.Error)
	assert.Nil(t, event.Data)
}

func TestEvent_Data_Success(t *testing.T) {
	matchMakerEvent := &MatchMakerEvent{
		Hash: common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
	}

	event := Event{Data: matchMakerEvent}

	assert.NoError(t, event.Error)
	assert.Equal(t, matchMakerEvent, event.Data)
}
