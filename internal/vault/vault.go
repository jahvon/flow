package vault

import (
	"crypto/sha512"
	"fmt"
	stdio "io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/flowexec/flow/internal/crypto"
	"github.com/flowexec/flow/internal/filesystem"
	"github.com/flowexec/flow/internal/logger"
)

const (
	EncryptionKeyEnvVar = "FLOW_VAULT_KEY"

	cacheDirName = "vault"
)

// Deprecated: Use the github.com/flowexec/vault package instead.
// This vault will be removed in a future release.
type Vault struct {
	cachedEncryptionKey string
	cachedData          *data
}

// Represents the data stored in the vault data file.
type data struct {
	LastUpdated string                 `yaml:"lastUpdated"`
	Secrets     map[string]SecretValue `yaml:"secrets"`
}

func RegisterEncryptionKey(key string) error {
	if err := filesystem.EnsureCachedDataDir(); err != nil {
		return err
	}

	fullPath := dataFilePath(key)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		return errors.New("encryption key already registered")
	}

	dir, _ := filepath.Split(fullPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return errors.Wrap(err, "unable to create vault data directory")
	}
	if _, err := os.Create(filepath.Clean(fullPath)); err != nil {
		return errors.Wrap(err, "unable to create vault data file")
	}

	return nil
}

func NewVault() *Vault {
	return &Vault{}
}

func (v *Vault) GetSecret(reference string) (SecretValue, error) {
	logger.Log().Debugf("getting secret with reference %s from vault", reference)
	d, err := v.loadData()
	if err != nil {
		return "", err
	} else if d == nil {
		return "", errors.New("no secrets found in vault")
	}

	secret, found := d.Secrets[reference]
	if !found {
		return "", fmt.Errorf("secret with reference %s not found", reference)
	}
	return secret, nil
}

func (v *Vault) GetAllSecrets() (map[string]SecretValue, error) {
	logger.Log().Debugf("getting all secrets from vault")
	d, err := v.loadData()
	if err != nil {
		return nil, err
	} else if d == nil {
		return nil, errors.New("no secrets found in vault")
	}
	return d.Secrets, nil
}

func (v *Vault) SetSecret(reference string, secret SecretValue) error {
	logger.Log().Debugf("setting secret with reference %s in vault", reference)
	if err := ValidateReference(reference); err != nil {
		return err
	}

	d, err := v.loadData()
	if err != nil {
		return err
	}

	if d.Secrets == nil {
		d.Secrets = make(map[string]SecretValue)
	}
	d.Secrets[reference] = secret

	return v.saveData(d)
}

func (v *Vault) DeleteSecret(reference string) error {
	logger.Log().Debugf("deleting secret with reference %s from vault", reference)
	d, err := v.loadData()
	if err != nil {
		return err
	}

	delete(d.Secrets, reference)

	return v.saveData(d)
}

func (v *Vault) RenameSecret(oldRef string, newRef string) error {
	logger.Log().Debugf("renaming secret with reference %s in vault", oldRef)

	d, err := v.loadData()
	if err != nil {
		return err
	}

	secret, exists := d.Secrets[oldRef]
	if !exists {
		return errors.Errorf("secret with reference %s does not exist in vault", oldRef)
	}

	_, exists = d.Secrets[newRef]
	if exists {
		return errors.Errorf("secret with reference %s already exists in vault", newRef)
	}

	if err := ValidateReference(newRef); err != nil {
		return err
	}

	d.Secrets[newRef] = secret

	delete(d.Secrets, oldRef)

	return v.saveData(d)
}

func (v *Vault) retrieveEncryptionKey() (string, error) {
	if v.cachedEncryptionKey != "" {
		return v.cachedEncryptionKey, nil
	}

	key, found := os.LookupEnv(EncryptionKeyEnvVar)
	if !found {
		return "", errors.New("encryption key not set")
	}
	if err := validateEncryptionKey(key); err != nil {
		return "", err
	}
	return key, nil
}

func (v *Vault) loadData() (*data, error) {
	if v.cachedData != nil {
		return v.cachedData, nil
	}

	key, err := v.retrieveEncryptionKey()
	if err != nil {
		return nil, err
	}

	fullPath := dataFilePath(key)
	file, err := os.Open(filepath.Clean(fullPath))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open vault data file")
	}
	defer file.Close()

	encryptedDataStr, err := stdio.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read vault data file")
	}

	if len(encryptedDataStr) == 0 {
		return &data{}, nil
	}

	dataStr, err := crypto.DecryptValue(key, string(encryptedDataStr))
	if err != nil {
		return nil, errors.Wrap(err, "unable to decrypt vault data")
	}

	var d data
	if err := yaml.Unmarshal([]byte(dataStr), &d); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal vault data")
	}
	return &d, nil
}

func (v *Vault) saveData(d *data) error {
	if d == nil {
		return nil
	}

	key, err := v.retrieveEncryptionKey()
	if err != nil {
		return err
	}

	d.LastUpdated = time.Now().Format(time.RFC3339)
	dataStr, err := yaml.Marshal(d)
	if err != nil {
		return errors.Wrap(err, "unable to marshal vault data")
	}
	encryptedDataStr, err := crypto.EncryptValue(key, string(dataStr))
	if err != nil {
		return errors.Wrap(err, "unable to encrypt vault data")
	}

	fullPath := dataFilePath(key)
	file, err := os.OpenFile(filepath.Clean(fullPath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to open vault data file")
	}
	defer file.Close()

	if _, err := file.WriteString(encryptedDataStr); err != nil {
		return errors.Wrap(err, "unable to write to vault data file")
	}
	return nil
}

func ValidateReference(reference string) error {
	if reference == "" {
		return errors.New("reference cannot be empty")
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !re.MatchString(reference) {
		return fmt.Errorf(
			"reference (%s) must only contain alphanumeric characters, dashes and/or underscores",
			reference,
		)
	}
	return nil
}

func validateEncryptionKey(key string) error {
	expectedDataPath := dataFilePath(key)
	if _, err := os.Stat(expectedDataPath); os.IsNotExist(err) {
		return errors.New("encryption key not recognized")
	}
	return nil
}

func dataFilePath(encryptionKey string) string {
	hasher := sha512.New()
	_, err := hasher.Write([]byte(encryptionKey))
	if err != nil {
		panic("unable to hash encryption key")
	}
	storageKey := crypto.EncodeValue(hasher.Sum(nil))
	return filepath.Join(filesystem.CachedDataDirPath(), cacheDirName, storageKey, "data")
}
