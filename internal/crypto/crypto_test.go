package crypto_test

import (
	"testing"

	"github.com/flowexec/flow/internal/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCrypto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crypto Suite")
}

var _ = Describe("GenerateKey", func() {
	It("generates a key", func() {
		key, err := crypto.GenerateKey()
		Expect(err).ToNot(HaveOccurred())
		Expect(key).ToNot(BeEmpty())

		decodedKey, err := crypto.DecodeValue(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedKey).ToNot(BeEmpty())
	})
})

var _ = Describe("DeriveKey", func() {
	It("derives a key from a password when salt is provided", func() {
		salt, err := crypto.GenerateKey()
		Expect(err).ToNot(HaveOccurred())
		decodedSalt, err := crypto.DecodeValue(salt)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedSalt).ToNot(BeEmpty())

		inputPassword := []byte("password")
		derivedKey, outSalt, err := crypto.DeriveKey(inputPassword, decodedSalt)
		Expect(err).ToNot(HaveOccurred())
		Expect(derivedKey).ToNot(BeEmpty())
		Expect(outSalt).To(Equal(salt))

		decodedDerivedKey, err := crypto.DecodeValue(derivedKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedDerivedKey).ToNot(BeEmpty())
	})

	It("derives a key from a password when the salt is not provided", func() {
		inputPassword := []byte("password")
		derivedKey, outSalt, err := crypto.DeriveKey(inputPassword, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(derivedKey).ToNot(BeEmpty())
		Expect(outSalt).ToNot(BeEmpty())

		decodedDerivedKey, err := crypto.DecodeValue(derivedKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedDerivedKey).ToNot(BeEmpty())
	})
})

var _ = Describe("EncryptValue and DecryptValue", func() {
	It("encrypts and decrypts a value", func() {
		masterKey, _ := crypto.GenerateKey()
		plaintext := "test value"
		encryptedValue, err := crypto.EncryptValue(masterKey, plaintext)
		Expect(err).ToNot(HaveOccurred())
		Expect(encryptedValue).ToNot(BeEmpty())
		Expect(encryptedValue).ToNot(Equal(plaintext))

		decryptedValue, err := crypto.DecryptValue(masterKey, encryptedValue)
		Expect(err).ToNot(HaveOccurred())
		Expect(decryptedValue).ToNot(BeEmpty())
		Expect(decryptedValue).To(Equal(plaintext))
	})
})
