// Package sse is the SSE Client for MEV-Share
package sse

// SSEClient is the SSE Client abstraction
type SSEClient interface {
	// Subscribe to events and returns a subscription
	Subscribe(eventChan chan<- Event) (SSESubscription, error)
	// MEV-Share event history
	EventHistoryInfo() (*EventHistoryInfo, error)
	// MEV-Share event history Params
	GetEventHistory(params EventHistoryParams) ([]EventHistory, error)
}

type SSESubscription interface {
	// To stop the subscription
	Stop()
}
