package backend

import (
	"errors"
	"time"

	"github.com/jahvon/flow/internal/backend/consts"
	"github.com/jahvon/flow/internal/crypto"
)

const defaultRememberMeDuration = 24 * time.Hour

type Config struct {
	Secret *SecretConfig `yaml:"secret"`
	Auth   *AuthConfig   `yaml:"auth"`
}

func NewConfig() *Config {
	return &Config{
		Secret: &SecretConfig{
			Backend: consts.EnvFileBackendName,
		},
		Auth: &AuthConfig{
			Backend:    consts.NoAuthBackendName,
			RememberMe: false,
		},
	}
}

func (c *Config) Validate() error {
	if c.Auth == nil {
		return errors.New("auth config not set")
	} else if c.Auth.Backend == "" {
		return errors.New("auth backend not set")
	}

	if c.Auth.PreferredMode != "" && c.Auth.PreferredMode != "password" && c.Auth.PreferredMode != "masterKey" {
		return errors.New("preferred mode must be either 'password' or 'masterKey'")
	}

	if c.Secret == nil {
		return errors.New("secret config not set")
	} else if c.Secret.Backend == "" {
		return errors.New("secret backend not set")
	}
	log.Trace().Msg("Backend config validated")
	return nil
}

type AuthConfig struct {
	Backend            consts.BackendName `yaml:"backend"`
	PreferredMode      consts.AuthMode    `yaml:"preferredMode"`
	RememberMe         bool               `yaml:"rememberMe"`
	RememberMeDuration time.Duration      `yaml:"rememberMeDuration"`

	curAuthBackend AuthBackend
}

func (a *AuthConfig) InitializeWithRandomMasterKey(backend AuthBackend) (string, error) {
	masterKey, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}
	var rememberMeDuration time.Duration
	if a.RememberMe {
		if a.RememberMeDuration == 0 {
			a.RememberMeDuration = defaultRememberMeDuration
		}
		rememberMeDuration = a.RememberMeDuration
	}
	a.curAuthBackend = backend
	return masterKey, a.curAuthBackend.SetNewMasterKey(masterKey, rememberMeDuration)
}

func (a *AuthConfig) LoginWithMasterKey(backend AuthBackend, masterKey string) error {
	if masterKey == "" {
		return errors.New("master key is empty")
	}

	var rememberMeDuration time.Duration
	if a.RememberMe {
		if a.RememberMeDuration == 0 {
			a.RememberMeDuration = defaultRememberMeDuration
		}
		rememberMeDuration = a.RememberMeDuration
	}
	a.curAuthBackend = backend
	return a.curAuthBackend.LoginWithMasterKey(masterKey, rememberMeDuration)
}

func (a *AuthConfig) InitializeWithPassword(backend AuthBackend, password string) (string, error) {
	var rememberMeDuration time.Duration
	if a.RememberMe {
		rememberMeDuration = a.RememberMeDuration
	}
	a.curAuthBackend = backend
	return a.curAuthBackend.SetNewPassword(password, rememberMeDuration)
}

func (a *AuthConfig) LoginWithPassword(backend AuthBackend, password string) error {
	if password == "" {
		return errors.New("password is empty")
	}

	var rememberMeDuration time.Duration
	if a.RememberMe {
		rememberMeDuration = a.RememberMeDuration
	}
	a.curAuthBackend = backend
	return a.curAuthBackend.LoginWithPassword(password, rememberMeDuration)
}

func (a *AuthConfig) CurrentBackend() AuthBackend {
	return a.curAuthBackend
}

type SecretConfig struct {
	Backend consts.BackendName `yaml:"backend"`

	curSecretBackend SecretBackend
}

func (s *SecretConfig) InitializeBackend() error {
	return s.curSecretBackend.InitializeBackend()
}

func (s *SecretConfig) CurrentBackend() SecretBackend {
	return s.curSecretBackend
}
