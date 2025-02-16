package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

// GenerateKey generates a random 32 byte key and returns it as a base64 encoded string.
func GenerateKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("error reading random bytes: %w", err)
	}
	return EncodeValue(key), nil
}

// DeriveKey derives a 32 byte key from the provided password and salt and returns
// the key and salt as base64 encoded strings.
// If salt is nil, a random salt will be generated.
func DeriveKey(password, salt []byte) (string, string, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return "", "", err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return "", "", err
	}

	return EncodeValue(key), EncodeValue(salt), nil
}

// EncodeValue encodes a byte slice as a base64 encoded string.
func EncodeValue(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// DecodeValue decodes a base64 encoded string into a byte slice.
func DecodeValue(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// EncryptValue encrypts a string using AES-256 and returns the encrypted value as a base64 encoded string.
// The encryption key used for encryption must be a base64 encoded string.
func EncryptValue(encryptionKey string, text string) (string, error) {
	decodedMasterKey, err := DecodeValue(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("error decoding master key: %w", err)
	}
	block, err := aes.NewCipher(decodedMasterKey)
	if err != nil {
		return "", fmt.Errorf("error creating new cipher: %w", err)
	}

	plaintext := []byte(text)
	// verify that the plaintext is not too long to fit in an int
	if len(plaintext) > 64*1024*1024 {
		return "", fmt.Errorf("plaintext too long to encrypt")
	}
	size := aes.BlockSize + len(plaintext)
	ciphertext := make([]byte, size)
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("error reading random bytes: %w", err)
	}

	cfb := cipher.NewCTR(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return string(ciphertext), nil
}

// DecryptValue decrypts a string using AES-256 and returns the decrypted value as a base64 encoded string.
// The master key used for decryption must be a base64 encoded string.
func DecryptValue(encryptionKey string, text string) (string, error) {
	decodedMasterKey, err := DecodeValue(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("error decoding master key: %w", err)
	}
	block, err := aes.NewCipher(decodedMasterKey)
	if err != nil {
		return "", fmt.Errorf("error creating new cipher: %w", err)
	}

	ciphertext := []byte(text)
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	plainText := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCTR(block, iv)
	cfb.XORKeyStream(plainText, ciphertext)
	return string(plainText), nil
}
