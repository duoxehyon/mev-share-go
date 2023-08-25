package sse

// SSE Client abstraction
type SSEApiClient interface {
	// Subscribe to events and returns a subscription
	Subscribe(eventChan chan<- Event) (SSESubscription, error)
	// MEV-Share event history
	EventHistoryInfo() (*EventHistoryInfo, error)
	// MEV-Share event history Params
	GetEventHistory(params EventHistoryParams) ([]EventHistory, error)
}

// Subscription type
type SSESubscription interface {
	// To stop the subscription
	Stop()
}
