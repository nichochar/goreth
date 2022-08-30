package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraConn = "https://mainnet.infura.io/v3/94112cb024c74fb697592b77c4819ff1"
	// Contract addresses
	bcsContractAddr  = "0xe182A80E76B1cF17D0eB018D563823357F1Ae296"
	nichocharEthAddr = "0x885F8588bB15a046f71bD5119f5BC3B67ee883d3"
	dwrAddr          = "0xD7029BDEa1c17493893AAfE29AAD69EF892B8ff2"
)

// CLI utility that prints a title nicely, directly to the console
// My Title becomes ->
// ########
// My Title
// ########
func printTitle(title string) {
	hashtags := strings.Repeat("#", len(title))
	fmt.Println("\n" + hashtags)
	fmt.Println(title)
	fmt.Println(hashtags + "\n")
}

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
	localFlag := flag.Bool("local", false, "a flag that, if passed, runs the script locally")
	fmt.Println("Fun with Go and Ethereum.")
	flag.Parse()
	fmt.Printf("localFlag: %v\n", localFlag)
	if *localFlag {
		mainLocal()
	} else {
		mainRemote()
	}
}
