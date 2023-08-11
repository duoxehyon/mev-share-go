package main

import (
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/sse"
)

func main() {
	client := sse.New("https://mev-share.flashbots.net")

	eventChan := make(chan sse.Event)
	sub, err := client.Subscribe(eventChan)

	if err != nil {
		log.Fatal(err)
	}

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
