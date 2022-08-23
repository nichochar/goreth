package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

const connectionString = "https://mainnet.infura.io/v3/94112cb024c74fb697592b77c4819ff1"

func main() {
	_, err := ethclient.Dial(connectionString)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("vim-go")
}
