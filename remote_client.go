package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func mainRemote() {
	printTitle("Client")
	infuraConn, err := infuraConn()
	fmt.Println(infuraConn)
	if err != nil {
		log.Fatal(err)
	}
	client, err := ethclient.Dial(infuraConn)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Chain ID: %v\n", chainID)
	suggestedGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Suggested Gas Price through the client with SuggestGasPrice(): %v wei\n", suggestedGasPrice)
	currentBlock, err := currentBlockNumber(client)
	if err != nil {
		log.Fatal(err)
	}
	currentBlockBigInt := new(big.Int).SetUint64(currentBlock)

	printTitle("Ledger Reads")
	for _, address := range []string{bcsContractAddr, nichocharEthAddr, dwrAddr} {
		ethBalance, err := getEthBalanceForAddr(client, address, nil)
		if err != nil {
			log.Fatal(err)
		}

		pendingEthBalance, err := getPendingEthBalanceForAddr(client, address)
		if err != nil {
			log.Fatal(err)
		}

		var contractIndicator = ""
		ok, err := isContractAddr(client, address)
		if err != nil {
			log.Fatal(err)
		}
		if ok {
			contractIndicator = "(smart contract)"
		}
		fmt.Printf("Account %v %v\n%20s: %.2f\n%20s: %.2f\n----\n", address[0:6], contractIndicator, "Balance", pendingEthBalance, "Pending Balance", ethBalance)

	}
	printTitle("Wallets")
	fmt.Println("Creating a wallet...")
	privateAddr, publicAddr, err := MakeWallet()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(privateAddr)
	spew.Dump(publicAddr)

	// Fun with blocks
	blockHeader, err := blockHeader(client, currentBlockBigInt)
	if err != nil {
		log.Fatal(err)
	}
	printTitle("Block Header")
	spew.Dump(blockHeader)

	block, err := blockByNumber(client, currentBlockBigInt)
	if err != nil {
		log.Fatal(err)
	}
	printTitle("Block")
	fmt.Printf("Current block number: %d\n", currentBlock)

	count, err := transactionCountInBlock(client, currentBlockBigInt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction count in block: %d (confirmed %d)\n", len(block.Transactions()), count)

	printTitle("Transactions")
	showCount := 3
	fmt.Printf("Showing %d transactions from block %d:\n", showCount, currentBlock)
	for idx, tx := range block.Transactions() {
		fmt.Printf("Hash().Hex(): %v\n", tx.Hash().Hex())
		fmt.Printf("Value(): %v\n", tx.Value().String())
		fmt.Printf("Gas(): %v\n", tx.Gas())
		fmt.Printf("GasPrice(): %v\n", tx.GasPrice().Uint64())
		fmt.Printf("Nonce(): %v\n", tx.Nonce())
		if len(tx.Data()) > 0 {
			fmt.Printf("Data(): %v....\n", tx.Data()[0:10])
		} else {
			fmt.Println("Empty call data for this tx.")
		}
		fmt.Printf("To().Hex(): %v\n", tx.To().Hex())
		// For getting the sender, we need the chainID (because of EIP-155)
		msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("msg.From().Hex(): %v\n", msg.From().Hex())

		// Check receipts
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Receipt.Status: %v\n", receipt.Status)
		fmt.Printf("Receipt.Logs: %v\n", receipt.Logs)
		fmt.Println()

		if idx > showCount {
			fmt.Printf("... and %d more tx not shown\n", len(block.Transactions())-showCount)
			break
		}
	}

}
