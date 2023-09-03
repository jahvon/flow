package defaults

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jahvon/flow/internal/backend"
	"github.com/jahvon/flow/internal/backend/consts"
)

type NoAuthBackend struct {
}

func AuthBackend() backend.AuthBackend {
	return &NoAuthBackend{}
}

func (n *NoAuthBackend) Name() consts.BackendName {
	return consts.NoAuthBackendName
}

func (n *NoAuthBackend) SetNewMasterKey(_ string, _ time.Duration) error {
	log.Trace().Msg("Skipping set of master key for no auth backend")
	return nil
}

func (n *NoAuthBackend) SetNewPassword(_ string, _ time.Duration) (string, error) {
	log.Trace().Msg("Skipping set of password for no auth backend")
	return "", nil
}

func (n *NoAuthBackend) LoginWithMasterKey(_ string, _ time.Duration) error {
	log.Trace().Msg("Skipping login with master key for no auth backend")
	return nil
}

func (n *NoAuthBackend) LoginWithPassword(_ string, _ time.Duration) error {
	log.Trace().Msg("Skipping login with password for no auth backend")
	return nil
}

func (n *NoAuthBackend) IsMasterKeyAuthorized() (bool, error) {
	log.Trace().Msg("Authorizing via no auth backend")
	return true, nil
}

func (n *NoAuthBackend) IsPasswordAuthorized(_ string) (bool, error) {
	log.Trace().Msg("Authorizing via no auth backend")
	return true, nil
}
