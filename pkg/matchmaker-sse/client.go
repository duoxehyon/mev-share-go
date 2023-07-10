package sse

// InternalClient is a client for the matchmaker.
type InternalClient struct {
	baseURL string // baseURL is the base URL for the matchmaker.
}

// NewMatchMakerSSE creates a new InternalClient for the matchmaker with the given base URL.
func NewMatchMakerSSE(baseURL string) *InternalClient {
	return &InternalClient{
		baseURL: baseURL,
	}
}
