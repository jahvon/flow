package keyring

import (
	"fmt"

	"github.com/zalando/go-keyring"

	"github.com/jahvon/tbox/internal/backend"
	"github.com/jahvon/tbox/internal/io"
)

const (
	// Keyring service name / prefix
	serviceName = "tbox"
)

var log = io.Log()

func getSecret(keyringContext, key string) (backend.Secret, error) {
	secret, err := keyring.Get(compositeServiceName(keyringContext), key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return backend.Secret(secret), nil
}

func setSecret(keyringContext, key string, secret backend.Secret) error {
	if err := keyring.Set(compositeServiceName(keyringContext), key, secret.String()); err != nil {
		return err
	}
	return nil
}

func deleteSecret(keyringContext string, key string) error {
	if err := keyring.Delete(compositeServiceName(keyringContext), key); err != nil {
		return err
	}
	return nil
}

func compositeServiceName(context string) string {
	if context != "" {
		return fmt.Sprintf("%s.%s", serviceName, context)
	}
	return serviceName
}
