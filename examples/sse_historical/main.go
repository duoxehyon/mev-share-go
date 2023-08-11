package main

import (
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/sse"
)

func main() {
	client := sse.New("https://mev-share.flashbots.net")

	info, err := client.EventHistoryInfo()
	if err != nil {
		log.Fatal(err)
	}

	query := sse.EventHistoryParams{
		BlockStart: info.MaxBlock - 100,
		BlockEnd:   info.MaxBlock,
		OffSet:     1,
	}

	history, err := client.GetEventHistory(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(history)
}
