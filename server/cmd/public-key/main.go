package main

import (
	"crypto/ecdsa"
	"os"

	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privateKeyStr := os.Args[1]
	if privateKeyStr == "" {
		panic("Please provide a private key.")
	}
	privateKeyBytes, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHex := hexutil.Encode(publicKeyBytes)
	fmt.Println(publicKeyHex)
	f, err := os.Create("public_key.hex")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(publicKeyHex)
	if err != nil {
		panic(err)
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)
	os.Exit(0)
}
