package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptBytes encrypts raw bytes using AES-256-GCM with PBKDF2 key derivation
// Returns the encrypted data as base64-encoded string (salt + ciphertext)
func EncryptBytes(password string, plainText []byte) (string, error) {
	// Generate random salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("salt generation err: %v", err)
	}

	// Derive a 32-byte key from the password (for AES-256)
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Create AES cipher with derived key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("cipher err: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("cipher GCM err: %v", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce err: %v", err)
	}

	// Encrypt the plaintext
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	// Combine salt + ciphertext
	finalData := append(salt, cipherText...)

	// Base64 encode the encrypted data
	encodedData := base64.StdEncoding.EncodeToString(finalData)

	return encodedData, nil
}
