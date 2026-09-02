package crypto

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	MasterKeyPrefix = "clam_sec_"
	KeyByteLength   = 32
	KeyHexLength    = 64
)

// GenerateMasterKey generates a fresh 256-bit cryptographically secure key
// with the standard 'clam_sec_' prefix.
func GenerateMasterKey() (string, []byte, error) {
	rawBytes := make([]byte, KeyByteLength)
	if _, err := rand.Read(rawBytes); err != nil {
		return "", nil, fmt.Errorf("failed to read random entropy: %w", err)
	}

	keyHex := hex.EncodeToString(rawBytes)
	keyStr := MasterKeyPrefix + keyHex
	return keyStr, rawBytes, nil
}

// ParseMasterKey validates and decodes a 'clam_sec_' prefixed 256-bit hex key.
func ParseMasterKey(keyStr string) ([]byte, error) {
	keyStr = strings.TrimSpace(keyStr)
	if !strings.HasPrefix(keyStr, MasterKeyPrefix) {
		return nil, errors.New("invalid master key format: missing 'clam_sec_' prefix")
	}

	hexPart := strings.TrimPrefix(keyStr, MasterKeyPrefix)
	if len(hexPart) != KeyHexLength {
		return nil, fmt.Errorf("invalid master key length: expected %d hex chars, got %d", KeyHexLength, len(hexPart))
	}

	rawBytes, err := hex.DecodeString(hexPart)
	if err != nil {
		return nil, fmt.Errorf("failed to decode master key hex: %w", err)
	}

	if len(rawBytes) != KeyByteLength {
		return nil, fmt.Errorf("invalid master key byte length: expected %d, got %d", KeyByteLength, len(rawBytes))
	}

	return rawBytes, nil
}

// EnsureMasterKey checks if an ENCRYPTION_KEY exists in the specified .env file.
// If missing or empty, it generates a fresh key, writes it to the .env file,
// and returns (keyStr, keyBytes, created=true, nil).
// If already present, it parses and returns the existing key (created=false).
func EnsureMasterKey(envPath string) (string, []byte, bool, error) {
	if envPath == "" {
		envPath = ".env"
	}

	// Check if file exists
	if _, err := os.Stat(envPath); err == nil {
		// File exists, attempt to read existing ENCRYPTION_KEY
		existingKey, err := readKeyFromEnv(envPath)
		if err == nil && existingKey != "" {
			rawBytes, err := ParseMasterKey(existingKey)
			if err == nil {
				return existingKey, rawBytes, false, nil
			}
		}
	}

	// Key is missing or file does not exist: generate a fresh one
	newKeyStr, rawBytes, err := GenerateMasterKey()
	if err != nil {
		return "", nil, false, err
	}

	if err := appendOrWriteEnv(envPath, newKeyStr); err != nil {
		return "", nil, false, fmt.Errorf("failed to write key to %s: %w", envPath, err)
	}

	return newKeyStr, rawBytes, true, nil
}

func readKeyFromEnv(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if strings.TrimSpace(parts[0]) == "ENCRYPTION_KEY" {
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			return val, nil
		}
	}

	return "", scanner.Err()
}

func appendOrWriteEnv(filePath string, keyStr string) error {
	// If file doesn't exist, create it with 0600 permissions
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		content := fmt.Sprintf("# Auto-generated clamav-service environment configuration\nENCRYPTION_KEY=%s\n", keyStr)
		return os.WriteFile(filePath, []byte(content), 0600)
	}

	// If file exists, check if ENCRYPTION_KEY line exists to replace or append
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ENCRYPTION_KEY=") {
			lines[i] = fmt.Sprintf("ENCRYPTION_KEY=%s", keyStr)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("ENCRYPTION_KEY=%s", keyStr))
	}

	newContent := strings.Join(lines, "\n")
	return os.WriteFile(filePath, []byte(newContent), 0600)
}
