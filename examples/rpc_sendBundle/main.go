package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/duoxehyon/mev-share-go/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	// Convert the signed raw tx hex string to bytes
	bytes, err := hex.DecodeString("02f86b0180843b9aca00852ecc889a0082520894c87037874aed04e51c29f582394217a0a2b89d808080c080a0a463985c616dd8ee17d7ef9112af4e6e06a27b071525b42182fe7b0b5c8b4925a00af5ca177ffef2ff28449292505d41be578bebb77110dfc09361d2fb56998260")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	txBytes := hexutil.Bytes(bytes)

	// Define the bundle transactions
	txns := []rpc.MevBundleBody{
		{
			Tx: &txBytes,
		},
	}

	inclusion := rpc.Inclusion{
		BlockNumber: 17891729,
	}

	// Make the bundle
	req := rpc.SendMevBundleArgs{
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
