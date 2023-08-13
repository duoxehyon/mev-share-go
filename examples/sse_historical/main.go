package main

import (
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/sse"
)

func main() {
	// Initialize the client
	client := sse.New("https://mev-share.flashbots.net")

	// Get info about mev-share transactions history
	info, err := client.EventHistoryInfo()
	if err != nil {
		log.Fatal(err)
	}

	// Make the query
	query := sse.EventHistoryParams{
		BlockStart: info.MaxBlock - 100,
		BlockEnd:   info.MaxBlock,
		OffSet:     1,
	}

	// Do query
	history, err := client.GetEventHistory(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(history)
}
