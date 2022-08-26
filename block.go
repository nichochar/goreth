package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func blockHeader(client *ethclient.Client, blockNumber *big.Int) (*types.Header, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func blockByNumber(client *ethclient.Client, blockNumber *big.Int) (*types.Block, error) {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func transactionCountInBlock(client *ethclient.Client, blockNumber *big.Int) (uint, error) {
	header, err := blockHeader(client, blockNumber)
	if err != nil {
		return 0, err
	}
	count, err := client.TransactionCount(context.Background(), header.Hash())
	if err != nil {
		return 0, nil
	}

	return count, nil
}
