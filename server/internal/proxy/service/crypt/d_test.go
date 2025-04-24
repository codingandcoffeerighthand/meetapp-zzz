package crypt_srv_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"golang.org/x/crypto/pbkdf2"
)

// EncryptedData khớp với định dạng của eth-crypto
type EncryptedData struct {
	Iv     string `json:"iv"`
	Ephemx string `json:"ephemPublicKey"`
	Ctext  string `json:"ciphertext"`
	Mac    string `json:"mac"`
}

func Test_a(t *testing.T) {
	// Private key Ethereum của bạn (giữ nó an toàn!)
	// Đây chỉ là khóa ví dụ, không bao giờ sử dụng nó trong sản phẩm thực tế!
	privateKeyHex := "70c40e1e7b43ead48cf4d66d04b420bae58f108d65d2a2b0efe7d8b820bf89ca" // Thay bằng private key thực

	// Giải mã chuỗi từ eth-crypto.cipher.stringify
	encryptedData, err := parseEncryptedData(
		`{"iv":"06ea2bb74bfd41ce6b3797c68c2589bf","ephemPublicKey":"04b22e1d99215ea1a06f945b624b79692d4a7cc87c503d5ce2abb866b70bfa0f04aeec1d2101f49c1989d3a262b3e8eb7863bc4413629a9dbb8849427ba9f06245","ciphertext":"c38dd3be07f2e6382c432af965955cebe7693f1d097297273c77da4a1e4f9fd2","mac":"ee2c467dde8958ea1c4e963e53c9e8be27b4ad827cca8ccf5f0e89c3b4398598"}`,
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Chuyển private key hex thành đối tượng ECDSA
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Giải mã dữ liệu
	decrypted, err := decrypt(privateKey, encryptedData)
	if err != nil {

		log.Println(err)
		t.Fail()
	}
	fmt.Println(decrypted)

}

// Phân tích chuỗi được mã hóa từ eth-crypto.cipher.stringify
func parseEncryptedData(encryptedString string) (d EncryptedData, e error) {
	e = json.Unmarshal([]byte(encryptedString), &d)
	return
}

// Giải mã dữ liệu bằng private key
func decrypt(privateKey *ecdsa.PrivateKey, encryptedData EncryptedData) (string, error) {
	// Chuyển đổi từ base64 sang bytes
	iv, err := hex.DecodeString(encryptedData.Iv)
	if err != nil {
		return "", fmt.Errorf("failed to decode iv: %w", err)
	}

	ephemPublicKeyStr, err := hex.DecodeString(encryptedData.Ephemx)
	if err != nil {
		return "", fmt.Errorf("failed to decode ephemeral public key: %w", err)
	}

	ciphertext, err := hex.DecodeString(encryptedData.Ctext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	mac, err := hex.DecodeString(encryptedData.Mac)
	if err != nil {
		return "", fmt.Errorf("failed to decode mac: %w", err)
	}

	// Chuyển đổi khóa công khai ephemeral từ dạng nén
	ephemeralPubKey, err := crypto.DecompressPubkey(ephemPublicKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse ephemeral public key: %w", err)
	}

	// Tính toán bí mật chia sẻ
	eciesPrivateKey := ecies.ImportECDSA(privateKey)
	eciesPublicKey := ecies.ImportECDSAPublic(ephemeralPubKey)

	sharedSecret, err := eciesPrivateKey.GenerateShared(eciesPublicKey, 16, 32)
	if err != nil {
		return "", fmt.Errorf("failed to generate shared secret: %w", err)
	}

	// Tạo khóa mã hóa và khóa MAC
	derivedKey := pbkdf2.Key(sharedSecret, nil, 1, 48, sha512.New)
	encryptionKey := derivedKey[:32]
	macKey := derivedKey[32:]

	// Kiểm tra MAC
	macData := make([]byte, 0, len(iv)+len(ephemPublicKeyStr)+len(ciphertext))
	macData = append(macData, iv...)
	macData = append(macData, ephemPublicKeyStr...)
	macData = append(macData, ciphertext...)

	calculatedMac := crypto.Keccak256(append(macKey, macData...))
	if !hmacEqual(calculatedMac, mac) {
		return "", fmt.Errorf("MAC verification failed")
	}

	// Giải mã AES-256-CBC
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Xử lý PKCS#7 padding
	paddingLen := int(plaintext[len(plaintext)-1])
	if paddingLen > len(plaintext) || paddingLen == 0 {
		return "", fmt.Errorf("invalid padding")
	}

	// Kiểm tra padding hợp lệ
	for i := len(plaintext) - paddingLen; i < len(plaintext); i++ {
		if plaintext[i] != byte(paddingLen) {
			return "", fmt.Errorf("invalid padding")
		}
	}

	// Loại bỏ padding
	plaintext = plaintext[:len(plaintext)-paddingLen]

	return string(plaintext), nil
}

// So sánh MAC theo thời gian không đổi
func hmacEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// PrintHex in ra dữ liệu dưới dạng hex để gỡ lỗi
func PrintHex(name string, data []byte) {
	fmt.Printf("%s: %s\n", name, hexutil.Encode(data))
}
