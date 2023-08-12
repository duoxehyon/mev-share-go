package shared

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ethereum/go-ethereum/common"
)

func TestLog_UnmarshalJSON(t *testing.T) {
	rawJSON := `{
		"address": "0x1234567890abcdef1234567890abcdef12345678",
		"topics": ["0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", "0x567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234"],
		"data": "0xdeadbeef"
	}`

	expectedAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedTopics := []common.Hash{
		common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
		common.HexToHash("0x567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234"),
	}
	expectedData, _ := hex.DecodeString("deadbeef")

	var log Log
	err := json.Unmarshal([]byte(rawJSON), &log)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, log.Address)
	assert.Equal(t, expectedTopics, log.Topics)
	assert.Equal(t, expectedData, log.Data)
}

func TestLog_UnmarshalJSON_NoData(t *testing.T) {
	rawJSON := `{
		"address": "0x1234567890abcdef1234567890abcdef12345678",
		"topics": ["0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"]
	}`

	expectedAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedTopics := []common.Hash{
		common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
	}

	var log Log
	err := json.Unmarshal([]byte(rawJSON), &log)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, log.Address)
	assert.Equal(t, expectedTopics, log.Topics)
	assert.Nil(t, log.Data)
}

func TestLog_UnmarshalJSON_EmptyData(t *testing.T) {
	rawJSON := `{
		"address": "0x1234567890abcdef1234567890abcdef12345678",
		"topics": ["0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"],
		"data": ""
	}`

	expectedAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedTopics := []common.Hash{
		common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
	}

	var log Log
	err := json.Unmarshal([]byte(rawJSON), &log)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, log.Address)
	assert.Equal(t, expectedTopics, log.Topics)
	assert.Nil(t, log.Data)
}
