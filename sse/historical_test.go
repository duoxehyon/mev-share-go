package sse

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInternalClient_EventHistoryInfo(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/history/info", r.URL.Path)
		info := EventHistoryInfo{
			Count:        100,
			MinBlock:     10000,
			MaxBlock:     20000,
			MinTimestamp: 1631419200,
			MaxLimit:     1000,
		}
		json.NewEncoder(w).Encode(info)
	}))
	defer mockServer.Close()

	client := &InternalClient{BaseURL: mockServer.URL}
	info, err := client.EventHistoryInfo()

	assert.NoError(t, err)
	assert.Equal(t, uint64(100), info.Count)
	assert.Equal(t, uint64(10000), info.MinBlock)
	assert.Equal(t, uint64(20000), info.MaxBlock)
	assert.Equal(t, uint64(1000), info.MaxLimit)
	assert.Equal(t, uint64(1631419200), info.MinTimestamp)
}

func TestInternalClient_GetEventHistory(t *testing.T) {
	params := EventHistoryParams{
		BlockStart:     10000,
		BlockEnd:       20000,
		TimestampStart: 1631419200,
		TimestampEnd:   1631422800,
		OffSet:         0,
	}

	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/history", r.URL.Path)

		var receivedParams EventHistoryParams
		err := json.NewDecoder(r.Body).Decode(&receivedParams)
		assert.NoError(t, err)
		assert.Equal(t, params, receivedParams)

		history := []EventHistory{
			{
				Block:     10001,
				Timestamp: 1631419260,
				Hint: MatchMakerEvent{
					Hash: common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				},
			},
			{
				Block:     10002,
				Timestamp: 1631419320,
				Hint: MatchMakerEvent{
					Hash: common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
				},
			},
		}
		json.NewEncoder(w).Encode(history)
	}))
	defer mockServer.Close()

	client := &InternalClient{BaseURL: mockServer.URL}
	history, err := client.GetEventHistory(params)

	assert.NoError(t, err)
	assert.Len(t, history, 2)

	assert.Equal(t, uint64(10001), history[0].Block)
	assert.Equal(t, uint64(1631419260), history[0].Timestamp)
	assert.Equal(t, common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"), history[0].Hint.Hash)

	assert.Equal(t, uint64(10002), history[1].Block)
	assert.Equal(t, uint64(1631419320), history[1].Timestamp)
	assert.Equal(t, common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"), history[1].Hint.Hash)
}
