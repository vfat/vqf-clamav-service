package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

// EncryptAESGCMBytes encrypts raw binary bytes using AES-256-GCM.
// It returns a slice containing [nonce 12-bytes + ciphertext + 16-byte auth tag].
func EncryptAESGCMBytes(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != KeyByteLength {
		return nil, fmt.Errorf("invalid key size: expected %d bytes for AES-256, got %d", KeyByteLength, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM AEAD: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate random nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// DecryptAESGCMBytes decrypts [nonce 12-bytes + ciphertext + 16-byte auth tag] binary payload using AES-256-GCM.
func DecryptAESGCMBytes(data []byte, key []byte) ([]byte, error) {
	if len(key) != KeyByteLength {
		return nil, fmt.Errorf("invalid key size: expected %d bytes for AES-256, got %d", KeyByteLength, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM AEAD: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (authentication tag mismatch or corrupted data): %w", err)
	}

	return plaintext, nil
}

// EncryptAESGCM encrypts a plaintext string using AES-256-GCM.
// It returns a standard base64-encoded string containing [nonce + ciphertext + auth_tag].
func EncryptAESGCM(plaintext string, key []byte) (string, error) {
	sealed, err := EncryptAESGCMBytes([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// DecryptAESGCM decrypts a base64-encoded [nonce + ciphertext + auth_tag] payload using AES-256-GCM.
func DecryptAESGCM(encodedCiphertext string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	plainBytes, err := DecryptAESGCMBytes(data, key)
	if err != nil {
		return "", err
	}

	return string(plainBytes), nil
}

// HashPassword creates a salted SHA-256 hash formatted as salt$hash.
func HashPassword(password string) string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)

	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)

	return fmt.Sprintf("%x$%x", salt, hash)
}

// VerifyPassword checks if a password matches a salt$hash string.
func VerifyPassword(password, storedHash string) bool {
	parts := strings.Split(storedHash, "$")
	if len(parts) != 2 {
		return false
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expectedHash, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	actualHash := hasher.Sum(nil)

	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1
}
