package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const keyBits = 2048

func generateRSAKey(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	// Ensure the private key is valid
	return privateKey, nil
}

func saveKey(privateKey *rsa.PrivateKey, filename string) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	file, err := os.Create(filename + ".pem")
	if err != nil {
		return fmt.Errorf("failed to create file private key: %w", err)
	}

	defer file.Close()

	err = pem.Encode(file, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to write data to file private key: %w", err)
	}
	publicKey := &privateKey.PublicKey
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyFileName := filename + ".pub.pem"
	publicKeyFile, err := os.Create(publicKeyFileName)
	if err != nil {
		return fmt.Errorf("failed to create file public key: %w", err)
	}
	defer publicKeyFile.Close()
	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to write data to file public key: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a filename to save the key.")
	}
	fileName := os.Args[1]
	if fileName == "" {
		panic("Filename cannot be empty.")
	}
	privateKey, err := generateRSAKey(keyBits)
	if err != nil {
		panic(err)
	}
	if err := saveKey(privateKey, fileName); err != nil {
		panic(err)
	}
	fmt.Println("Successfully generated RSA key pair.")
}
