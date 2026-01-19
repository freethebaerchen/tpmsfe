package encryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func DecryptFileContent(password string, filePath string) error {
	// Read base64 encoded file
	encodedData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file err: %v", err)
	}

	// Remove all whitespace and line breaks
	cleanedData := cleanBase64(encodedData)

	// Decode from base64
	encryptedData, err := base64.StdEncoding.DecodeString(cleanedData)
	if err != nil {
		return fmt.Errorf("base64 decode err: %v", err)
	}

	// Extract salt (first 32 bytes)
	if len(encryptedData) < 32 {
		return fmt.Errorf("invalid encrypted file format")
	}
	salt := encryptedData[:32]
	cipherText := encryptedData[32:]

	// Derive the same key from password
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher err: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("cipher GCM err: %v", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Decrypt
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return fmt.Errorf("decrypt err: %v", err)
	}

	// Write decrypted content
	err = os.WriteFile(filePath, plainText, 0644)
	if err != nil {
		return fmt.Errorf("write file err: %v", err)
	}

	return nil
}

func cleanBase64(data []byte) string {
	// Remove common whitespace characters
	cleaned := bytes.ReplaceAll(data, []byte("\n"), []byte(""))
	cleaned = bytes.ReplaceAll(cleaned, []byte("\r"), []byte(""))
	cleaned = bytes.ReplaceAll(cleaned, []byte(" "), []byte(""))
	cleaned = bytes.ReplaceAll(cleaned, []byte("\t"), []byte(""))

	return string(cleaned)
}
