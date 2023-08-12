package main

import (
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	// Make the bundle
	req := rpc.MevSendBundleParams{
		Body:      txns,
		Inclusion: inclusion,
	}

	// Add overrides
	overrides := rpc.SimBundleOverrides{
		ParentBlock: hexutil.Uint64(1000000),
	}

	// Simulate the bundle
	res, err := client.SimBundle(req, overrides)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
