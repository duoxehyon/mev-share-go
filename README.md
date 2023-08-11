# MEV-Share Client

A simple and clean client library for MEV-Share written in Golang.

Based on [MEV-Share Spec](https://github.com/flashbots/mev-share).

# Usage

To add library to your project:

``go get github.com/duoxehyon/mev-share-go``

## Subscribing to MEV-Share Events

Simplistic api for subscribing to MEV-Share events 

```go
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
	}
}

```

## Sending bundles 

Example on how to send bundles using this client

```go
package main

import (
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/rpc"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fbSigningKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		log.Fatal(err)
	}

	client := rpc.NewClient("https://relay.flashbots.net", fbSigningKey)

	// Define the bundle transactions
	txns := []rpc.BundleItem{
		rpc.MevShareTxHash{
			Hash: "0x......", // hash from an mev-share event
		},
		rpc.SignedRawTx{
			Tx:        "0x......", // signed raw transaction
			CanRevert: false,
		},
	}

	inclusion := rpc.Inclusion{
		Block: 17891729,
	}

	req := rpc.MevSendBundleParams{
		Body:      txns,
		Inclusion: inclusion,
	}

	res, err := client.SendBundle(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.BundleHash)
}
```

For more usage examples, explore the /examples directory in the library repository.

## License

Licensed under:

* MIT license ([LICENSE-MIT](LICENSE-MIT) or
  https://opensource.org/licenses/MIT)
