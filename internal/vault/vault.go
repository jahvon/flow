package vault

import (
	"crypto/sha512"
	"errors"
	"fmt"
	stdio "io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/crypto"
	"github.com/jahvon/flow/internal/io"
)

const (
	EncryptionKeyEnvVar = "FLOW_VAULT_KEY"

	cacheDirName = "vault"
)

var log = io.Log().With().Str("pkg", "vault").Logger()

type Vault struct {
	cachedEncryptionKey string
	cachedData          *data
}

// Represents the data stored in the vault data file.
type data struct {
	LastUpdated string            `yaml:"lastUpdated"`
	Secrets     map[string]Secret `yaml:"secrets"`
}

func RegisterEncryptionKey(key string) error {
	log.Trace().Msg("registering encryption key")
	if err := file.EnsureCachedDataDir(); err != nil {
		return err
	}

	fullPath := dataFilePath(key)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		return fmt.Errorf("encryption key already registered")
	}

	dir, _ := filepath.Split(fullPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("unable to create vault directory - %w", err)
	}
	if _, err := os.Create(filepath.Clean(fullPath)); err != nil {
		return fmt.Errorf("unable to create vault data file - %w", err)
	}

	return nil
}

func NewVault() *Vault {
	return &Vault{}
}

func (v *Vault) GetSecret(reference string) (Secret, error) {
	log.Trace().Msgf("getting secret with reference %s", reference)
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

func (v *Vault) GetAllSecrets() (map[string]Secret, error) {
	log.Trace().Msg("getting all secrets")
	d, err := v.loadData()
	if err != nil {
		return nil, err
	} else if d == nil {
		return nil, errors.New("no secrets found in vault")
	}
	return d.Secrets, nil
}

func (v *Vault) SetSecret(reference string, secret Secret) error {
	log.Trace().Msgf("setting secret with reference %s", reference)
	if err := ValidateReference(reference); err != nil {
		return err
	}

	d, err := v.loadData()
	if err != nil {
		return err
	}

	if d.Secrets == nil {
		d.Secrets = make(map[string]Secret)
	}
	d.Secrets[reference] = secret

	return v.saveData(d)
}

func (v *Vault) DeleteSecret(reference string) error {
	log.Trace().Msgf("deleting secret with reference %s", reference)
	d, err := v.loadData()
	if err != nil {
		return err
	}

	delete(d.Secrets, reference)

	return v.saveData(d)
}

func (v *Vault) retrieveEncryptionKey() (string, error) {
	if v.cachedEncryptionKey != "" {
		log.Debug().Msg("using cached encryption key")
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
		log.Debug().Msg("using cached vault data")
		return v.cachedData, nil
	}
	log.Trace().Msg("loading vault data from file")

	key, err := v.retrieveEncryptionKey()
	if err != nil {
		return nil, err
	}

	fullPath := dataFilePath(key)
	file, err := os.Open(filepath.Clean(fullPath))
	if err != nil {
		return nil, fmt.Errorf("unable to open vault data file - %w", err)
	}
	defer file.Close()

	encryptedDataStr, err := stdio.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read vault data file - %w", err)
	}

	if len(encryptedDataStr) == 0 {
		return &data{}, nil
	}

	dataStr, err := crypto.DecryptValue(key, string(encryptedDataStr))
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt vault data - %w", err)
	}

	var d data
	if err := yaml.Unmarshal([]byte(dataStr), &d); err != nil {
		return nil, fmt.Errorf("unable to unmarshal vault data - %w", err)
	}
	return &d, nil
}

func (v *Vault) saveData(d *data) error {
	if d == nil {
		return nil
	}
	log.Trace().Msg("saving vault data to file")

	key, err := v.retrieveEncryptionKey()
	if err != nil {
		return err
	}

	d.LastUpdated = time.Now().Format(time.RFC3339)
	dataStr, err := yaml.Marshal(d)
	if err != nil {
		return fmt.Errorf("unable to marshal vault data - %w", err)
	}
	encryptedDataStr, err := crypto.EncryptValue(key, string(dataStr))
	if err != nil {
		return fmt.Errorf("unable to encrypt vault data - %w", err)
	}

	fullPath := dataFilePath(key)
	file, err := os.OpenFile(filepath.Clean(fullPath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open vault data file - %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(encryptedDataStr); err != nil {
		return fmt.Errorf("unable to write to vault data file - %w", err)
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
		log.Panic().Err(err).Msg("unable to hash encryption key")
	}
	storageKey := crypto.EncodeValue(hasher.Sum(nil))
	return filepath.Join(file.CachedDataDirPath(), cacheDirName, storageKey, "data")
}
