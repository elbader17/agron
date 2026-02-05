package agron

import (
	"bytes"
	"testing"
)

func generateTestKey() []byte {
	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		key[i] = byte(i)
	}
	return key
}

func TestNewVault_ValidKey(t *testing.T) {
	key := generateTestKey()
	vault, err := NewVault(key)
	if err != nil {
		t.Fatalf("NewVault failed with valid key: %v", err)
	}
	if vault == nil {
		t.Fatal("NewVault returned nil vault")
	}
}

func TestNewVault_InvalidKeySize(t *testing.T) {
	invalidKeys := [][]byte{
		make([]byte, 16),
		make([]byte, 24),
		make([]byte, 31),
		make([]byte, 33),
	}

	for _, key := range invalidKeys {
		_, err := NewVault(key)
		if err == nil {
			t.Errorf("NewVault should fail with %d-byte key", len(key))
		}
		if err != ErrInvalidKeySize {
			t.Errorf("Expected ErrInvalidKeySize, got: %v", err)
		}
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := generateTestKey()
	vault, err := NewVault(key)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	plaintext := []byte("Hello, this is a secret message!")
	context := []byte("user-context-123")

	ciphertext, err := vault.Encrypt(plaintext, context)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("Ciphertext should not equal plaintext")
	}

	decrypted, err := vault.Decrypt(ciphertext, context)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted text doesn't match. Got %s, want %s", decrypted, plaintext)
	}
}

func TestDecrypt_WrongContext(t *testing.T) {
	key := generateTestKey()
	vault, err := NewVault(key)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	plaintext := []byte("Secret data")
	correctContext := []byte("correct-context")
	wrongContext := []byte("wrong-context")

	ciphertext, err := vault.Encrypt(plaintext, correctContext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = vault.Decrypt(ciphertext, wrongContext)
	if err == nil {
		t.Error("Decrypt should fail with wrong context")
	}
	if err != ErrDecryption {
		t.Errorf("Expected ErrDecryption, got: %v", err)
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key := generateTestKey()
	vault, err := NewVault(key)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	plaintext := []byte("Important data")
	context := []byte("context")

	ciphertext, err := vault.Encrypt(plaintext, context)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)-1] ^= 0xFF

	_, err = vault.Decrypt(tampered, context)
	if err == nil {
		t.Error("Decrypt should fail with tampered ciphertext")
	}
	if err != ErrDecryption {
		t.Errorf("Expected ErrDecryption, got: %v", err)
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	key := generateTestKey()
	vault, err := NewVault(key)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}

	shortData := make([]byte, 4)
	_, err = vault.Decrypt(shortData, []byte("context"))
	if err == nil {
		t.Error("Decrypt should fail with too-short ciphertext")
	}
	if err != ErrDecryption {
		t.Errorf("Expected ErrDecryption, got: %v", err)
	}
}
