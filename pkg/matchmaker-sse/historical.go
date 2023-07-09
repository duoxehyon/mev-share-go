package sse

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
	// The timestamp when the event was emitted.
	Timestamp uint64 `json:"timestamp,omitempty"`
	// Mev-share tx hint
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
