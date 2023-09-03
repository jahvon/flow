package backend

import (
	"time"

	"github.com/jahvon/tbox/internal/backend/consts"
	"github.com/jahvon/tbox/internal/io"
)

var log = io.Log()

type AuthBackend interface {
	Name() consts.BackendName
	SetNewMasterKey(masterKey string, rememberMeDuration time.Duration) error
	SetNewPassword(password string, rememberMeDuration time.Duration) (string, error)
	LoginWithMasterKey(masterKey string, rememberMeDuration time.Duration) error
	LoginWithPassword(password string, rememberMeDuration time.Duration) error
	IsMasterKeyAuthorized() (bool, error)
	IsPasswordAuthorized(password string) (bool, error)
}

type SecretBackend interface {
	Name() consts.BackendName
	InitializeBackend() error
	GetSecret(context, key string) (Secret, error)
	SetSecret(context, key string, secret Secret) error
	DeleteSecret(context, key string) error
}
