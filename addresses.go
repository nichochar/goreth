package main

import (
	"context"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func isValidEthAddr(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)

}

func isContractAddr(client *ethclient.Client, addrString string) (bool, error) {
	address := common.HexToAddress(addrString)
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is the latest block
	if err != nil {
		return false, err
	}
	return len(bytecode) > 0, nil
}
