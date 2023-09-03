package keyring

import (
	"errors"
	"fmt"
	"time"

	"github.com/jahvon/flow/internal/backend"
	"github.com/jahvon/flow/internal/backend/consts"
	"github.com/jahvon/flow/internal/crypto"
)

const (
	authBackendContext      = "kr_auth"
	authCheckSecretKey      = "check"
	authCheckSecretData     = "Z3VpUGg1U01WNzg1bjd0QWRjWWdUSWM4RjdhSnRVZGU=" // encoded
	authMasterKeySecretKey  = "masterKey"
	authExpirationSecretKey = "expiration"
	authSaltSecretKey       = "salt"
)

type authBackend struct {
	masterKey string
}

func AuthBackend() backend.AuthBackend {
	return &authBackend{}
}

func (a *authBackend) Name() consts.BackendName {
	return consts.KeyringBackendName
}

// SetNewMasterKey sets the master key for the keyring auth backend.
// If rememberMeDuration is greater than 0, the master key will be persisted
// for the specified duration.
func (a *authBackend) SetNewMasterKey(masterKey string, rememberMeDuration time.Duration) error {
	log.Trace().Msg("Setting keyring auth master key")
	a.masterKey = masterKey
	if a.masterKey == "" {
		return errors.New("master key cannot be empty")
	} else if len(masterKey) < 32 {
		return errors.New("master key must be at least 32 characters")
	}

	if err := setCheckSecret(masterKey); err != nil {
		return fmt.Errorf("failed to set check secret in keyring: %w", err)
	}
	if rememberMeDuration > 0 {
		if err := a.persistMasterKey(rememberMeDuration); err != nil {
			return fmt.Errorf("failed to persist master key: %w", err)
		}
	}

	return nil
}

// LoginWithMasterKey logs in with the master key for the keyring auth backend.
// If authorized and rememberMeDuration is greater than 0, the
// expiration will be updated.
func (a *authBackend) LoginWithMasterKey(masterKey string, rememberMeDuration time.Duration) error {
	log.Trace().Msg("Logging in with keyring auth master key")
	a.masterKey = masterKey
	if authorized, err := a.IsMasterKeyAuthorized(); err != nil {
		return fmt.Errorf("unable to check if master key is authorized: %w", err)
	} else if !authorized {
		return errors.New("master key not authorized")
	} else if rememberMeDuration > 0 {
		log.Debug().Msg("Master key authorized; updating expiration")
		if err := setExpiration(time.Now().Add(rememberMeDuration)); err != nil {
			return fmt.Errorf("unable to update master key expiration: %w", err)
		}
	}
	return nil
}

// IsMasterKeyAuthorized checks if the master key is authorized.
func (a *authBackend) IsMasterKeyAuthorized() (bool, error) {
	masterKey := a.masterKey
	if masterKey == "" {
		persistedMk, err := a.getPersistedMasterKey()
		if err != nil {
			log.Err(err).Msg("Failed to get persisted master key")
		}
		masterKey = persistedMk
	}

	if masterKey == "" {
		return false, nil
	}

	authorized, err := compareCheckSecret(masterKey)
	if err != nil {
		return false, fmt.Errorf("unable to compare check secret: %w", err)
	}

	return authorized, nil
}

// SetNewPassword sets the password for the keyring auth backend.
// If rememberMeDuration is greater than 0, the password will be persisted
// for the specified duration.
// The current derived master key will be returned along with any errors.
func (a *authBackend) SetNewPassword(password string, rememberMeDuration time.Duration) (string, error) {
	log.Trace().Msg("Setting keyring auth password")

	if set, err := a.passwordSet(); err != nil {
		log.Error().Err(err).Msg("Unable to check if password is set")
	} else if set {
		log.Warn().Msg("Updating an existing password")
	}

	masterKey, salt, err := crypto.DeriveKey([]byte(password), nil)
	if err != nil {
		return "", fmt.Errorf("unable to derive master key from password: %w", err)
	}

	if err := setSecret(authBackendContext, authSaltSecretKey, backend.Secret(salt)); err != nil {
		return "", fmt.Errorf("unable to set salt secret: %w", err)
	}

	if err := a.SetNewMasterKey(masterKey, rememberMeDuration); err != nil {
		return "", fmt.Errorf("unable to initialize auth backend with master key: %w", err)
	}

	return masterKey, nil
}

// LoginWithPassword logs in with the password for the keyring auth backend.
// If authorized and rememberMeDuration is greater than 0, the
// expiration will be updated.
func (a *authBackend) LoginWithPassword(password string, rememberMeDuration time.Duration) error {
	log.Trace().Msg("Logging in with keyring auth password")
	if authorized, err := a.IsPasswordAuthorized(password); err != nil {
		return fmt.Errorf("unable to check if password is authorized: %w", err)
	} else if !authorized {
		return errors.New("password not authorized")
	} else if rememberMeDuration > 0 {
		log.Debug().Msg("Password authorized; updating expiration")
		if err := setExpiration(time.Now().Add(rememberMeDuration)); err != nil {
			return fmt.Errorf("unable to update password expiration: %w", err)
		}
	}
	return nil
}

// IsPasswordAuthorized checks if the password is authorized.
func (a *authBackend) IsPasswordAuthorized(password string) (bool, error) {
	if set, err := a.passwordSet(); err != nil {
		return false, fmt.Errorf("unable to check if password is set: %w", err)
	} else if !set {
		return false, nil
	}

	salt, err := getSecret(authBackendContext, authSaltSecretKey)
	if err != nil {
		return false, fmt.Errorf("unable to get salt secret: %w", err)
	}

	decodedSalt, err := crypto.DecodeValue(salt.String())
	if err != nil {
		return false, fmt.Errorf("unable to decode salt: %w", err)
	}

	masterKey, _, err := crypto.DeriveKey([]byte(password), decodedSalt)
	if err != nil {
		return false, fmt.Errorf("unable to derive master key from password: %w", err)
	}

	authorized, err := compareCheckSecret(masterKey)
	if err != nil {
		return false, fmt.Errorf("unable to compare check secret: %w", err)
	}

	a.masterKey = masterKey
	return authorized, nil
}

func (a *authBackend) persistMasterKey(duration time.Duration) error {
	if a.masterKey == "" {
		return errors.New("master key not set")
	}

	if duration <= 0 {
		duration = 24 * time.Hour
	}
	log.Trace().Str("duration", duration.String()).Msg("Persisting keyring auth master key")

	if err := setSecret(authBackendContext, authMasterKeySecretKey, backend.Secret(a.masterKey)); err != nil {
		return fmt.Errorf("failed to set master key secret in keyring: %w", err)
	}
	if err := setExpiration(time.Now().Add(duration)); err != nil {
		return fmt.Errorf("failed to set master key expiration data in keyring: %w", err)
	}
	return nil
}

func (a *authBackend) getPersistedMasterKey() (string, error) {
	log.Trace().Msg("Getting persisted keyring auth master key")

	expiration, err := getSecret(authBackendContext, authExpirationSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to get master key expiration from keyring: %w", err)
	} else if expiration.Empty() {
		return "", errors.New("master key expiration not set")
	}
	expirationTime, err := time.Parse(time.RFC3339, expiration.String())
	if err != nil {
		return "", fmt.Errorf("failed to parse master key expiration: %w", err)
	}
	// If the expiration time is in the past, delete the master key and return an error.
	if time.Now().After(expirationTime) {
		err := deleteSecret(authBackendContext, authMasterKeySecretKey)
		if err != nil {
			return "", fmt.Errorf("failed to delete master key secret after expiration: %w", err)
		}
		return "", errors.New("master key expired")
	}

	storedMasterKey, err := getSecret(authBackendContext, authMasterKeySecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to get master key from keyring: %w", err)
	} else if storedMasterKey.Empty() {
		return "", errors.New("keyring auth backend not initialized")
	}

	decryptedMasterKey, err := storedMasterKey.Decrypt(a.masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt master key: %w", err)
	}

	return decryptedMasterKey.String(), nil
}

func (a *authBackend) passwordSet() (bool, error) {
	salt, err := getSecret(authBackendContext, authSaltSecretKey)
	if err != nil {
		return false, fmt.Errorf("unable to get salt secret: %w", err)
	}

	return !salt.Empty(), nil
}

func setExpiration(expiration time.Time) error {
	if err := setSecret(authBackendContext, authExpirationSecretKey, backend.Secret(expiration.Format(time.RFC3339))); err != nil {
		return fmt.Errorf("failed to set master key expiration data in keyring: %w", err)
	}
	return nil
}

func setCheckSecret(masterKey string) error {
	encryptedSecret, err := checkSecretPlaintext().Encrypt(masterKey)
	if err != nil {
		return err
	}
	if err := setSecret(authBackendContext, authCheckSecretKey, encryptedSecret); err != nil {
		return err
	}
	log.Trace().Msg("Check secret updated")
	return nil
}

func getCheckSecret(masterKey string) (backend.Secret, error) {
	encryptedSecret, err := getSecret(authBackendContext, authCheckSecretKey)
	if err != nil {
		return "", err
	} else if encryptedSecret.Empty() {
		return "", errors.New("keyring auth backend missing check secret")
	}

	decryptedSecret, err := encryptedSecret.Decrypt(masterKey)
	if err != nil {
		return "", err
	}

	return decryptedSecret, nil
}

func compareCheckSecret(masterKey string) (bool, error) {
	decryptedStoredSecret, err := getCheckSecret(masterKey)
	if err != nil {
		return false, err
	}
	return decryptedStoredSecret.String() == checkSecretPlaintext().String(), nil
}

func checkSecretPlaintext() backend.Secret {
	decodedCheckSecret, _ := crypto.DecodeValue(authCheckSecretData)
	return backend.Secret(decodedCheckSecret)
}
