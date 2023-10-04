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

	// Signed raw transaction in hex
	txn := "0x02f86b0180843b9aca00852ecc889a0282520894c87037874aed04e51c29f582394217a0a2b89d808080c080a0a463985c616dd8ee17d7ef9112af4e6e06a27b071525b42182fe7b0b5c8b4925a00af5ca177ffef2ff28449292505d41be578bebb77110dfc09361d2fb56998260" // signed raw transaction

	options := rpc.PrivateTxOptions{
		Hints: rpc.Hints{
			CallData:        true,
			ContractAddress: true,
			Logs:            true,
		},
		MaxBlockNumber: hexutil.Uint64(100),
		Builders:       []string{"flashbots"},
	}

	// Send the private transaction
	res, err := client.SendPrivateTransaction(txn, &options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.String())
}
