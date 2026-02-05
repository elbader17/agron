// Package agron implements AES-GCM authenticated encryption with pluggable key loaders.
package agron

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	ErrInvalidKeySize = errors.New("agron: key must be 32 bytes (hex encoded string usually)")
	ErrDecryption     = errors.New("agron: decryption failed (integrity check or context mismatch)")
)

type Vault struct {
	gcm cipher.AEAD
}

// NewVault recibe la llave en BYTES crudos (ya decodificados)
func NewVault(masterKey []byte) (*Vault, error) {
	if len(masterKey) != 32 {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Vault{gcm: gcm}, nil
}

// Encrypt: Encripta data usando un contexto (AAD)
func (v *Vault) Encrypt(plaintext, context []byte) ([]byte, error) {
	nonce := make([]byte, v.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return v.gcm.Seal(nonce, nonce, plaintext, context), nil
}

// Decrypt: Desencripta validando el contexto
func (v *Vault) Decrypt(ciphertext, context []byte) ([]byte, error) {
	nonceSize := v.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrDecryption
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := v.gcm.Open(nil, nonce, encryptedData, context)
	if err != nil {
		return nil, ErrDecryption
	}
	return plaintext, nil
}
