package keyring

import (
	"github.com/jahvon/tbox/internal/backend"
	"github.com/jahvon/tbox/internal/backend/consts"
)

type secretBackend struct {
}

func SecretBackend() backend.SecretBackend {
	return &secretBackend{}
}

func (s *secretBackend) Name() consts.BackendName {
	return consts.KeyringBackendName
}

func (s *secretBackend) InitializeBackend() error {
	return nil
}

func (s *secretBackend) GetSecret(keyringContext, key string) (backend.Secret, error) {
	log.Trace().Msgf("getting secret for keyring context %s and key %s", keyringContext, key)
	return getSecret(keyringContext, key)
}

func (s *secretBackend) SetSecret(keyringContext, key string, secret backend.Secret) error {
	log.Trace().Msgf("updating secret for keyring context %s and key %s", keyringContext, key)
	return setSecret(keyringContext, key, secret)
}

func (s *secretBackend) DeleteSecret(keyringContext, key string) error {
	log.Trace().Msgf("deleting secret for keyring context %s and key %s", keyringContext, key)
	return deleteSecret(keyringContext, key)
}
