package domain

import (
	"encoding/base64"
)

type ns struct{}

func NewDumDecryptService() (*ns, error) {
	return &ns{}, nil
}

func (ns *ns) Decrypt(data string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	return string(decodeBytes), err
}

func (ns *ns) Encrypt(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
