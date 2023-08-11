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

	txn := "0x......" // signed raw transaction

	options := rpc.PrivateTxOptions{
		Hints: rpc.Hints{
			CallData:        true,
			ContractAddress: true,
			Logs:            true,
		},
	}

	res, err := client.SendPrivateTransaction(txn, &options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.String())
}
