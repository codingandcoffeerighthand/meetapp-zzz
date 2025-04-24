package crypt_srv

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"proxy-srv/internal/proxy/configs"
)

type srv struct {
	privateKey *rsa.PrivateKey
}

func NewCryptService(cfg configs.Config) (*srv, error) {
	pemBytes, err := os.ReadFile(cfg.PrivateKeyRSA)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key %s", cfg.PrivateKeyRSA)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &srv{
		privateKey: privateKey,
	}, nil
}

func (s *srv) Decrypt(data string) (string, error) {
	hash := crypto.SHA256
	label := []byte("")
	decryptedData, err := rsa.DecryptOAEP(hash.New(), rand.Reader, s.privateKey, []byte(data), label)
	return string(decryptedData), err
}
