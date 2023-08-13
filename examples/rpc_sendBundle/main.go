package main

import (
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/rpc"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Flashbots header signing key
	fbSigningKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the client
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

	// Make the bundle
	req := rpc.MevSendBundleParams{
		Body:      txns,
		Inclusion: inclusion,
	}

	// Send bundle
	res, err := client.SendBundle(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.BundleHash.String())
}
