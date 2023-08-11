# MEV-Share Client

Client library for MEV-Share written in Golang.

Based on [MEV-Share Spec](https://github.com/flashbots/mev-share).


## Subscribing to MEV-Share Events

Simplistic api for subscribing to MEV-Share events 

```go
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
	}
}

```


## License

Licensed under:

* MIT license ([LICENSE-MIT](LICENSE-MIT) or
  https://opensource.org/licenses/MIT)
