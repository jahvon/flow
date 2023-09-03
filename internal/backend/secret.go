package backend

import (
	"strings"

	"github.com/jahvon/tbox/internal/crypto"
)

const (
	encryptionPrefix = "tbox:v1:"
)

type Secret string

func (s Secret) String() string {
	return string(s)
}

func (s Secret) Empty() bool {
	return string(s) == ""
}

func (s Secret) IsEncrypted() bool {
	return strings.HasPrefix(s.String(), encryptionPrefix)
}

func (s Secret) Decrypt(masterKey string) (Secret, error) {
	if s.IsEncrypted() {
		canDecrypt := masterKey != ""
		if canDecrypt {
			log.Trace().Msg("Decrypting a secret")
			decryptedValue, err := crypto.DecryptValue(masterKey, s.withoutPrefix().String())
			if err != nil {
				return "", err
			}
			return Secret(decryptedValue), nil
		} else {
			log.Debug().Msg("masterKey not set; skipping secret decryption")
		}
	}
	return s, nil
}

func (s Secret) Encrypt(masterKey string) (Secret, error) {
	if !s.IsEncrypted() {
		canEncrypt := masterKey != ""
		if canEncrypt {
			log.Trace().Msg("Encrypting a secret")
			encryptedValue, err := crypto.EncryptValue(masterKey, s.String())
			if err != nil {
				return "", err
			}
			return Secret(encryptedValue).withPrefix(), nil
		} else {
			log.Debug().Msg("masterKey not set; skipping secret encryption")
		}
	}
	return s, nil
}

func (s Secret) withoutPrefix() Secret {
	return Secret(strings.TrimPrefix(s.String(), encryptionPrefix))
}

func (s Secret) withPrefix() Secret {
	return Secret(encryptionPrefix + s.String())
}
