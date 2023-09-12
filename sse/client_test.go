package sse

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		_, err := w.Write([]byte("data: {\"some\":\"event\"}\n\n"))
		if err != nil {
			panic(err)
		}
	}))
}

func TestInternalClient_Subscribe(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	client := New(server.URL)

	eventChan := make(chan Event, 1)
	_, err := client.Subscribe(eventChan)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	select {
	case event := <-eventChan:
		if event.Data == nil {
			t.Error("Expected event data, got nil")
		}
	case <-time.After(5 * time.Second):
		t.Error("Timed out waiting for event")
	}
}

func TestSubscription_Stop(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	client := New(server.URL)

	eventChan := make(chan Event, 1)
	subscription, err := client.Subscribe(eventChan)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	subscription.Stop()

	_, ok := <-eventChan

	assert.Equal(t, false, ok)
}
