# MEV-Share Client

A simple and clean client library for MEV-Share written in Golang.

Based on [MEV-Share Spec](https://github.com/flashbots/mev-share).

# Usage

To add library to your project:

``go get github.com/duoxehyon/mev-share-go``

## Subscribing to MEV-Share Events

Begin by subscribing to MEV-Share events with the following code snippet:
import the sse client
```go
import (
	"github.com/duoxehyon/mev-share-go/sse" // Import the SSE client

	"fmt"
	"log"
)

func main() {
	// Create a new client
	client := sse.New("https://mev-share.flashbots.net")

	// Make event channel for receiving events
	eventChan := make(chan sse.Event)

	// Subscribe to events
	sub, err := client.Subscribe(eventChan)
	if err != nil {
		log.Fatal(err)
	}

	// Read events
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
	"github.com/duoxehyon/mev-share-go/rpc" // Import the RPC client

	"fmt"
	"log"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Initialize flashbots auth key. this key is only used for searcher reputation and is not required to hold any funds   
	fbSigningKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new client
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

	// Define the block to include / max block to include
	inclusion := rpc.Inclusion{
		Block: 17891729,
	}

	// Make the bundle
	req := rpc.MevSendBundleParams{
		Body:      txns,
		Inclusion: inclusion,
	}

	// Send it
	res, err := client.SendBundle(req)
	if err != nil {
		log.Fatal(err)
	}

	// Returns a bundle hash which is unique to every bundle 
	fmt.Println(res.BundleHash)
}
```

## Sending Private Transactions

```go
func main() {
	// auth key
	fbSigningKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		log.Fatal(err)
	}

	// init client
	client := rpc.NewClient("https://relay.flashbots.net", fbSigningKey)

	// signed raw transaction
	txn := "0x......" 

	// Extra params while sending private tx 
	options := rpc.PrivateTxOptions{
		Hints: rpc.Hints{
			CallData:        true,
			ContractAddress: true,
			Logs:            true,
		},
	}

	// Send tx
	res, err := client.SendPrivateTransaction(txn, &options)
	if err != nil {
		log.Fatal(err)
	}

	// Print tx hash
	fmt.Println(res.String())
}
```

For more usage examples, explore the /examples directory in the library repository.

## License

Licensed under:

* MIT license ([LICENSE-MIT](LICENSE-MIT) or
  https://opensource.org/licenses/MIT)
