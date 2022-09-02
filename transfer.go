package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const localChainID = 1337

func transferEthLocal(
	client *ethclient.Client,
	to string,
	privateKeyHex string, // this is also how we get the sender public address
	weiAmount *big.Int) (*types.Transaction, error) {

	fromAddress, err := addressFromPrivateKeyString(account0.privateKey)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v -> (%veth) -> %v\n", fromAddress, weiToEth(weiAmount), to)

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		return nil, err
	}

	// We are opinionated here because we know
	// the exact gasLimit for an eth transfer
	gasLimit := uint64(21000)

	// This is another opinionated decision, we're using the recommended
	// value (with no buffer). We could adopt different strategies based
	// on greediness vs certainty the tx will go through.
	suggestedGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	toAddress := common.HexToAddress(account1.address)

	// The big one! Data field is nil since we're just sending eth
	tx := types.NewTransaction(nonce, toAddress, weiAmount, gasLimit, suggestedGasPrice, nil)

	// Now that we've created the transaction, we need to sign it:
	// for this we use the privateKey from the sender (account0)
	//chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	privateKey, err := privateKeyFromHex(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// 1337 is a magical number because the chain we're running on is
	// a private geth implementation(ganache)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(localChainID)), privateKey)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Success, the tx got mined: %v\n", signedTx.Hash().Hex())
	return signedTx, nil
}
