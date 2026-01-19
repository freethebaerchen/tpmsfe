package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func EncryptFileContent(password string, filePath string) error {
	// Read plaintext file
	plainText, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file err: %v", err)
	}

	// Generate random salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("salt generation err: %v", err)
	}

	// Derive a 32-byte key from the password (for AES-256)
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Create AES cipher with derived key
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher err: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("cipher GCM err: %v", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce err: %v", err)
	}

	// Encrypt the plaintext
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	// Combine salt + ciphertext
	finalData := append(salt, cipherText...)

	// Base64 encode the encrypted data
	encodedData := base64.StdEncoding.EncodeToString(finalData)

	// Format with line breaks every 64 characters
	formattedData := formatBase64(encodedData)

	// Write formatted content to file
	err = os.WriteFile(filePath, []byte(formattedData), 0644)
	if err != nil {
		return fmt.Errorf("write file err: %v", err)
	}

	return nil
}

func formatBase64(encoded string) string {
	var formatted strings.Builder

	for i := 0; i < len(encoded); i += 64 {
		end := i + 64
		if end > len(encoded) {
			end = len(encoded)
		}
		formatted.WriteString(encoded[i:end])
		if end < len(encoded) {
			formatted.WriteString("\n")
		}
	}

	return formatted.String()
}
