package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// EncryptAESGCM encrypts a plaintext string using AES-256-GCM.
// It returns a standard base64-encoded string containing [nonce + ciphertext + auth_tag].
func EncryptAESGCM(plaintext string, key []byte) (string, error) {
	if len(key) != KeyByteLength {
		return "", fmt.Errorf("invalid key size: expected %d bytes for AES-256, got %d", KeyByteLength, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM AEAD: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate random nonce: %w", err)
	}

	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// DecryptAESGCM decrypts a base64-encoded [nonce + ciphertext + auth_tag] payload using AES-256-GCM.
func DecryptAESGCM(encodedCiphertext string, key []byte) (string, error) {
	if len(key) != KeyByteLength {
		return "", fmt.Errorf("invalid key size: expected %d bytes for AES-256, got %d", KeyByteLength, len(key))
	}

	data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM AEAD: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext payload too short: missing valid nonce")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed (authentication tag mismatch or corrupted data): %w", err)
	}

	return string(plaintext), nil
}
