package encryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// DecryptBytes decrypts base64-encoded encrypted data using AES-256-GCM with PBKDF2 key derivation
// Returns the decrypted plaintext bytes
func DecryptBytes(password string, encryptedBase64 string) ([]byte, error) {
	// Remove all whitespace and line breaks (defensive measure in case base64 has whitespace)
	cleaned := bytes.ReplaceAll([]byte(encryptedBase64), []byte("\n"), []byte(""))
	cleaned = bytes.ReplaceAll(cleaned, []byte("\r"), []byte(""))
	cleaned = bytes.ReplaceAll(cleaned, []byte(" "), []byte(""))
	cleaned = bytes.ReplaceAll(cleaned, []byte("\t"), []byte(""))
	cleanedData := string(cleaned)

	// Decode from base64
	encryptedData, err := base64.StdEncoding.DecodeString(cleanedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode err: %v", err)
	}

	// Extract salt (first 32 bytes)
	if len(encryptedData) < 32 {
		return nil, fmt.Errorf("invalid encrypted data format")
	}
	salt := encryptedData[:32]
	cipherText := encryptedData[32:]

	// Derive the same key from password
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher err: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher GCM err: %v", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Decrypt
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt err: %v", err)
	}

	return plainText, nil
}
