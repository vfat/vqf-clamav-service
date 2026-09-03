package crypto

import (
	"testing"
)

func TestAESGCM_Roundtrip(t *testing.T) {
	_, keyBytes, err := GenerateMasterKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	testCases := []string{
		"super_secret_webhook_url_12345",
		"bot_token_789123:AAHzxy98412",
		"{\"smtp_user\":\"admin\",\"smtp_pass\":\"P@ssw0rd123!\"}",
		"",
	}

	for _, original := range testCases {
		cipherText, err := EncryptAESGCM(original, keyBytes)
		if err != nil {
			t.Fatalf("encryption failed for input '%s': %v", original, err)
		}

		if original != "" && cipherText == original {
			t.Errorf("ciphertext must not equal plaintext")
		}

		decrypted, err := DecryptAESGCM(cipherText, keyBytes)
		if err != nil {
			t.Fatalf("decryption failed for input '%s': %v", original, err)
		}

		if decrypted != original {
			t.Errorf("expected decrypted text '%s', got '%s'", original, decrypted)
		}
	}
}

func TestAESGCM_RandomNonce(t *testing.T) {
	_, keyBytes, _ := GenerateMasterKey()
	input := "same_secret_string"

	cipher1, err := EncryptAESGCM(input, keyBytes)
	if err != nil {
		t.Fatalf("encrypt 1 failed: %v", err)
	}

	cipher2, err := EncryptAESGCM(input, keyBytes)
	if err != nil {
		t.Fatalf("encrypt 2 failed: %v", err)
	}

	if cipher1 == cipher2 {
		t.Errorf("expected different ciphertexts due to random 96-bit nonce, but got identical strings")
	}
}

func TestAESGCM_InvalidKeyLength(t *testing.T) {
	shortKey := []byte("too_short_key_16b")
	_, err := EncryptAESGCM("secret", shortKey)
	if err == nil {
		t.Errorf("expected error with 16-byte key (needs 32-byte AES-256), got nil")
	}

	_, err = DecryptAESGCM("invalid_cipher", shortKey)
	if err == nil {
		t.Errorf("expected error with 16-byte key on decrypt, got nil")
	}
}

func TestAESGCM_TamperedCiphertext(t *testing.T) {
	_, keyBytes, _ := GenerateMasterKey()
	cipherText, err := EncryptAESGCM("sensitive_data", keyBytes)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Ubah 1 karakter pada ciphertext (merusak tag otentikasi GCM)
	tampered := []byte(cipherText)
	if tampered[len(tampered)-1] == 'A' {
		tampered[len(tampered)-1] = 'B'
	} else {
		tampered[len(tampered)-1] = 'A'
	}

	_, err = DecryptAESGCM(string(tampered), keyBytes)
	if err == nil {
		t.Errorf("expected authentication tag failure on tampered ciphertext, got nil")
	}
}

func TestPasswordHashing(t *testing.T) {
	pwd := "123456"
	hash := HashPassword(pwd)

	if !VerifyPassword(pwd, hash) {
		t.Errorf("expected password '%s' to verify against hash '%s'", pwd, hash)
	}

	if VerifyPassword("wrongpassword", hash) {
		t.Errorf("expected wrong password to fail verification")
	}
}
