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
		return "", fmt.Errorf("error reading random bytes: %s", err)
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
// The master key used for encryption must be a base64 encoded string.
func EncryptValue(masterKey string, text string) (string, error) {
	decodedMasterKey, err := DecodeValue(masterKey)
	if err != nil {
		return "", fmt.Errorf("error decoding master key: %s", err)
	}
	block, err := aes.NewCipher(decodedMasterKey)
	if err != nil {
		return "", fmt.Errorf("error creating new cipher: %s", err)
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("error reading random bytes: %s", err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return fmt.Sprintf("%s", ciphertext), nil
}

// DecryptValue decrypts a string using AES-256 and returns the decrypted value as a base64 encoded string.
// The master key used for decryption must be a base64 encoded string.
func DecryptValue(masterKey string, text string) (string, error) {
	decodedMasterKey, err := DecodeValue(masterKey)
	if err != nil {
		return "", fmt.Errorf("error decoding master key: %s", err)
	}
	block, err := aes.NewCipher(decodedMasterKey)
	if err != nil {
		return "", fmt.Errorf("error creating new cipher: %s", err)
	}

	ciphertext := []byte(text)
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	plainText := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plainText, ciphertext)
	return fmt.Sprintf("%s", plainText), nil
}
