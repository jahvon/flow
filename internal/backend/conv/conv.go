package conv

import (
	"fmt"

	"github.com/jahvon/tbox/internal/backend"
	"github.com/jahvon/tbox/internal/backend/consts"
	"github.com/jahvon/tbox/internal/backend/defaults"
	"github.com/jahvon/tbox/internal/backend/keyring"
)

func AuthBackendFromName(name consts.BackendName) (backend.AuthBackend, error) {
	var authBackend backend.AuthBackend
	switch name {
	case consts.NoAuthBackendName:
		authBackend = defaults.AuthBackend()
	case consts.KeyringBackendName:
		authBackend = keyring.AuthBackend()
	default:
		return nil, fmt.Errorf("unknown auth backend - %v", name)
	}
	return authBackend, nil
}

func SecretBackendFromName(name consts.BackendName) (backend.SecretBackend, error) {
	var secretBackend backend.SecretBackend
	switch name {
	case consts.EnvFileBackendName:
		secretBackend = defaults.SecretBackend()
	case consts.KeyringBackendName:
		secretBackend = keyring.SecretBackend()
	default:
		return nil, fmt.Errorf("unknown secret backend - %v", name)
	}
	return secretBackend, nil
}
