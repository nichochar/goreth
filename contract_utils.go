package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func buildAuthForTx(client *ethclient.Client, privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	// Get the nonce for this address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECSDA (unexpected)")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	// Get the gas price
	suggestedGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(localChainID))
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = suggestedGasPrice
	return auth, nil
}

func deployContract(client *ethclient.Client, privateKey *ecdsa.PrivateKey) (*Storage, error) {
	auth, err := buildAuthForTx(client, privateKey)
	if err != nil {
		return nil, err
	}

	address, tx, instance, err := DeployStorage(auth, client)
	fmt.Printf("Address of contract: %v\n", address)
	fmt.Printf("Tx: %v\n: ", tx)

	return instance, nil
}

func loadContract(client *ethclient.Client, contractAddressHex string) (*Storage, error) {
	address := common.HexToAddress(contractAddressHex)
	instance, err := NewStorage(address, client)
	return instance, err
}
