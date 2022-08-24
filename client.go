package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraConn = "https://mainnet.infura.io/v3/94112cb024c74fb697592b77c4819ff1"
	// Assumes you're running the following in another process
	localGanacheConn = "http://localhost:8545"

	// Contract addresses
	bcsContractAddr  = "0xe182A80E76B1cF17D0eB018D563823357F1Ae296"
	nichocharEthAddr = "0x885F8588bB15a046f71bD5119f5BC3B67ee883d3"
	pushixEthAddr    = "0x1dA4FDf7029bDf8ff11f28141a659f6563940642"
)

// 1 * wei * 10^18 = 1 eth
func weiToEth(weiBalance *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(weiBalance.String())
	ethBalance := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethBalance
}

func getEthBalanceForAddr(client *ethclient.Client, addr string, block *big.Int) (*big.Float, error) {
	account := common.HexToAddress(addr)
	// nil is block number, passing that gives the latest balance
	balance, err := client.BalanceAt(context.Background(), account, block)
	if err != nil {
		return big.NewFloat(0), nil
	}
	return weiToEth(balance), nil

}

func currentBlockNumber(client *ethclient.Client) (uint64, error) {
	return client.BlockNumber(context.Background())
}

func getPendingEthBalanceForAddr(client *ethclient.Client, addr string) (*big.Float, error) {
	account := common.HexToAddress(addr)
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	if err != nil {
		return big.NewFloat(0), nil
	}
	return weiToEth(pendingBalance), nil
}

func main() {
	fmt.Println("Fun with Go and Ethereum.")
	client, err := ethclient.Dial(infuraConn)
	if err != nil {
		log.Fatal(err)
	}

	currentBlock, err := currentBlockNumber(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current block number: %d\n", currentBlock)

	for _, address := range []string{bcsContractAddr, nichocharEthAddr, pushixEthAddr} {
		ethBalance, err := getEthBalanceForAddr(client, address, nil)
		if err != nil {
			log.Fatal(err)
		}

		pendingEthBalance, err := getPendingEthBalanceForAddr(client, address)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Account %v\n%20s: %.2f\n%20s: %.2f\n----\n", address[0:6], "Balance", pendingEthBalance, "Pending Balance", ethBalance)

	}
}