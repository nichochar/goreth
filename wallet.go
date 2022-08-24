package main

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func MakeWallet() (pub string, priv string, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateAddr := hexutil.Encode(privateKeyBytes)
	fmt.Printf("New private key bytes:%v\n", privateAddr)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	publicAddr := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("Address: %v\n", publicAddr)

	// The public address is simply the Keccak-256 hash of the public key, and then we take the last 40 characters (20 bytes) and prefix it with 0x
	// Let's reconstruct it manually with lower level crypto calls
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Printf("Regenerated address: %v\n", hexutil.Encode(hash.Sum(nil)[12:]))

	if !isValidEthAddr(publicAddr) {
		return "", "", fmt.Errorf("The address created was invalid, %s", publicAddr)
	}

	return publicAddr, privateAddr, nil
}
