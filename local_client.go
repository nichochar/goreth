package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	// You can always re-deploy it with `go run . -local -deploy`
	// and update this value if you wish.
	// Current owner: account0
	contractAddressHex = "0x0Cf777784463a2c643c762959985690E6b0609E6"
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

func printBalancesForAccounts1And2(client *ethclient.Client) {
	account0Balance, err := getEthBalanceForAddr(client, account0.address, nil)
	if err != nil {
		log.Fatal(err)
	}
	account1Balance, err := getEthBalanceForAddr(client, account1.address, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("---")
	fmt.Printf("Account0 balance: %veth\n", account0Balance)
	fmt.Printf("Account1 balance: %veth\n", account1Balance)
	fmt.Println("---")
}

func mainLocal(deploy bool, writeNum int, readNum bool) {
	fmt.Println("Running against local client (ganache)...")
	client, err := ethclient.Dial(localGanacheConn)
	if err != nil {
		log.Fatal(err)
	}
	printBalancesForAccounts1And2(client)

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	transferEthLocal(client, account1.address, account0.privateKey, value)
	printBalancesForAccounts1And2(client)

	if deploy {
		fmt.Println("Deploying contract...")
		privateKeyAccount0, err := privateKeyFromHex(account0.privateKey)
		if err != nil {
			log.Fatal(err)
		}
		storage, err := deployContract(client, privateKeyAccount0)
		if err != nil {
			log.Print("Error deploying contract")
			log.Fatal(err)
		}
		fmt.Printf("We have storage: %v\n", storage)
	}

	// Let's interact with the contract
	instance, err := loadContract(client, contractAddressHex)
	privateKey, err := privateKeyFromHex(account1.privateKey)
	if err != nil {
		log.Fatal(err)
	}

	if writeNum != 0 {
		auth, err := buildAuthForTx(client, privateKey)
		if err != nil {
			log.Fatal(err)
		}
		// Let's write to the contract with account1!

		auth.Value = big.NewInt(0)     // in wei
		auth.GasLimit = uint64(300000) // in units

		num := big.NewInt(int64(writeNum))
		tx, err := instance.Store(auth, num)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Transaction went through! We saved %d to contract store\n", num)
		fmt.Printf("(tx: %v)\n", tx.Hash().Hex())
	}

	if readNum {
		result, err := instance.Retrieve(&bind.CallOpts{})
		if err != nil {
			log.Print("Error reading from store")
			log.Fatal(err)
		}
		fmt.Printf("Read value from store: %d\n", result)

	}

}
