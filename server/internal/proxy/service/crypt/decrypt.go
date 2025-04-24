package crypt_srv

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"proxy-srv/internal/proxy/configs"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

type s struct {
	pk string
}

func NewDecryptService(cfg *configs.Web3Config) *s {
	return &s{
		pk: cfg.PrivateKey,
	}
}

type EncryptedDataEthCrypto struct {
	IV             string `json:"iv"`
	EphemPublicKey string `json:"ephemPublicKey"`
	Ciphertext     string `json:"ciphertext"`
	MAC            string `json:"mac"`
}

func (s *s) Decrypt(data string) (string, error) {
	privateKeyHex := s.pk

	// Example encrypted data (hex-encoded)
	e := &EncryptedDataEthCrypto{}
	err := json.Unmarshal([]byte(data), e)

	// Convert private key from hex to *ecdsa.PrivateKey
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Printf("Error decoding private key: %v\n", err)
		return "", err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		fmt.Printf("Error converting to ECDSA private key: %v\n", err)
		return "", err
	}

	// Get the public address for verification
	publicKeyECDSA := privateKey.Public().(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHex := hexutil.Encode(publicKeyBytes)
	fmt.Println(publicKeyHex)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Ethereum address: %s\n", address.Hex())

	// Convert encrypted data from hex to bytes

	// Convert the ECDSA private key to an ECIES private key
	eciesPrivateKey := ecies.ImportECDSA(privateKey)

	i, err := hex.DecodeString(e.IV)
	if err != nil {
		return "", err
	}

	r, err := hex.DecodeString(e.EphemPublicKey)
	if err != nil {
		return "", err
	}
	c, err := hex.DecodeString(e.Ciphertext)
	if err != nil {
		return "", err

	}
	t, err := hex.DecodeString(e.MAC)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	buffer.Write(i)
	buffer.Write(r)
	buffer.Write(c)
	buffer.Write(t)

	// Decrypt the data
	decryptedData, err := eciesPrivateKey.Decrypt(buffer.Bytes(), nil, nil)
	if err != nil {
		fmt.Printf("Error decrypting data: %v\n", err)
		return "", err
	}

	return string(decryptedData), nil
}

type ns struct{}

func NewDumDecryptService() (*ns, error) {
	return &ns{}, nil
}

func (ns *ns) Decrypt(data string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	return string(decodeBytes), err
}
