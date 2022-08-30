package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Account struct {
	address    string
	privateKey string
	balanceEth float64
}

const (
	// Assumes you're running the following in another process
	localGanacheConn = "http://localhost:8545"
)

var account0 = Account{
	address:    "0x1B470Cb74266D8F305D1C7f1fdbFeE60b1dc8f31",
	privateKey: "0x9299af2b6158dffada96d4a31a2582dafb469bca7f71dfd3bd7af12b77432198",
	balanceEth: 1000.,
}

var account1 = Account{
	address:    "0xe2A2be2D1c765F57fD9531391Cf2049A64AbA0D6",
	privateKey: "0x7ecfe554aaa0ae991cda90d481751d8443c2cff60dc709b37172a9e010003ded",
	balanceEth: 1000.,
}

func mainLocal() {
	fmt.Println("Running against local client (ganache)...")
	client, err := ethclient.Dial(localGanacheConn)
	if err != nil {
		log.Fatal(err)
	}

	account0Balance, err := getEthBalanceForAddr(client, account0.address, nil)
	if err != nil {
		log.Fatal(err)
	}
	account1Balance, err := getEthBalanceForAddr(client, account1.address, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Account0 balance: %v\n", account0Balance)
	fmt.Printf("Account1 balance: %v\n", account1Balance)

	fromAddress, err := addressFromPrivateKeyString(account0.privateKey)
	if err != nil {
		log.Fatal(err)
	}
	if account0.address != fromAddress {
		log.Fatal("account0.address should be the same as fromAddress")
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)

	gasLimit := uint64(21000) // in units

	// Default Gas Price in ganache: 2000000000
	gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	suggestedGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Gas price: %v\nSuggested: %v\n", gasPrice, suggestedGasPrice)

	toAddress := common.HexToAddress(account1.address)

	// The big one! Data field is nil since we're just sending eth
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Now that we've created the transaction, we need to sign it:
	// for this we use the privateKey from the sender (account0)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := privateKeyFromHex(account0.privateKey)
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0x77006fcb3938f648e2cc65bafd27dec30b9bfbe9df41f78498b9c8b7322a249e

}
