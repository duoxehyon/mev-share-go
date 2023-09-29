package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/flashbots/mev-share-node/mevshare"
)

func main() {
	// Flashbots header signing key
	fbSigningKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the client
	client := rpc.NewClient("https://relay.flashbots.net", fbSigningKey)

	// Convert the hex string to bytes
	bytes, err := hex.DecodeString("signed raw transaction")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	txBytes := hexutil.Bytes(bytes)

	// Define the bundle transactions
	txns := []mevshare.MevBundleBody{
		{
			Tx: &txBytes,
		},
	}

	inclusion := mevshare.MevBundleInclusion{
		BlockNumber: 17891729,
	}

	// Make the bundle
	req := mevshare.SendMevBundleArgs{
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
