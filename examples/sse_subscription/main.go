package main

import (
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/sse"
)

func main() {
	// Initialize client
	client := sse.New("https://mev-share.flashbots.net")

	// Channel to send messages to
	eventChan := make(chan sse.Event)
	// Subscribe to events from mev-share node
	sub, err := client.Subscribe(eventChan)
	if err != nil {
		log.Fatal(err)
	}

	// Listen for events
	for {
		event := <-eventChan
		if event.Error != nil {
			fmt.Println("Error occured: ", event.Error)
		}

		fmt.Println(event)

		sub.Stop()
		break
	}
}
