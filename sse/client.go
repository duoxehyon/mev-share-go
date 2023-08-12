package sse

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
)

// InternalClient is a client for the matchmaker
type InternalClient struct {
	BaseURL string // BaseURL is the base URL for the matchmaker
}

// New creates a new InternalClient for the matchmaker with the given base URL
func New(baseURL string) *InternalClient {
	return &InternalClient{
		BaseURL: baseURL,
	}
}

// Subscription represents a subscription to matchmaker events
type Subscription struct {
	client    http.Client
	stopper   chan struct{}
	scanner   *bufio.Scanner
	eventChan chan<- Event
}

// Subscribe to matchmaker events and returns a type that can be used to control the subscription
func (c *InternalClient) Subscribe(eventChan chan<- Event) (*Subscription, error) {
	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	sub := &Subscription{
		client:    client,
		eventChan: eventChan,
		stopper:   make(chan struct{}),
		scanner:   bufio.NewScanner(resp.Body),
	}

	go sub.readEvents()

	return sub, nil
}

// readEvents reads the events and sends them to the event channel
func (s *Subscription) readEvents() {
	for s.scanner.Scan() {
		data := s.scanner.Text()
		if data == ":ping" || data == "" {
			continue
		}

		data = strings.TrimPrefix(data, "data: ")

		var event MatchMakerEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			s.eventChan <- Event{
				Error: err,
			}
			continue
		}

		select {
		case <-s.stopper:
			return
		default:
			s.eventChan <- Event{
				Data: &event,
			}
		}
	}
}

// Stop stops the subscription to matchmaker events
func (s *Subscription) Stop() {
	s.stopper <- struct{}{}
}
