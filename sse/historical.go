package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// For querying historical mev-share transactions
type EventHistoryParams struct {
	BlockStart     uint64 `json:"blockStart,omitempty"`
	BlockEnd       uint64 `json:"blockEnd:,omitempty"`
	TimestampStart uint64 `json:"timestampStart:,omitempty"`
	TimestampEnd   uint64 `json:"timestampEnd:,omitempty"`
	OffSet         uint64 `json:"offset:,omitempty"`
}

// Single historical mev-share transaction
type EventHistory struct {
	// Block number of event's block
	Block uint64 `json:"block,omitempty"`
	// The timestamp when the event was emitted
	Timestamp uint64 `json:"timestamp,omitempty"`
	// Mev-share tx hint.
	Hint MatchMakerEvent `json:"hint,omitempty"`
}

// Info on mev-share historical data
type EventHistoryInfo struct {
	Count        uint64 `json:"count"`
	MinBlock     uint64 `json:"minBlock"`
	MaxBlock     uint64 `json:"maxBlock"`
	MinTimestamp uint64 `json:"minTimestamp"`
	MaxLimit     uint64 `json:"maxLimit"`
}

// Gets info about historical mev-share data
func (c *InternalClient) EventHistoryInfo() (*EventHistoryInfo, error) {
	url := c.BaseURL + "/api/v1/history/info"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var eventHistoryInfo EventHistoryInfo
	err = json.NewDecoder(resp.Body).Decode(&eventHistoryInfo)
	if err != nil {
		return nil, err
	}

	return &eventHistoryInfo, nil
}

// Gets historical mev-share data
func (c *InternalClient) GetEventHistory(params EventHistoryParams) ([]EventHistory, error) {
	url := c.BaseURL + "/api/v1/history"

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var eventHistory []EventHistory
	err = json.NewDecoder(resp.Body).Decode(&eventHistory)
	if err != nil {
		return nil, err
	}

	return eventHistory, nil
}
