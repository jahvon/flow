package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flowexec/tuikit/io"
	"github.com/flowexec/vault"

	"github.com/flowexec/flow/internal/filesystem"
	"github.com/flowexec/flow/internal/utils"
)

const (
	DefaultVaultKeyEnv      = "FLOW_VAULT_KEY"
	DefaultVaultIdentityEnv = "FLOW_VAULT_IDENTITY"
	LegacyVaultReservedName = "legacy"

	v2CacheDataDir = "vaults"
)

type Vault = vault.Provider
type VaultConfig = vault.Config

func NewAES256Vault(logger io.Logger, name, storagePath, keyEnv, keyFile, logLevel string) {
	if keyEnv == "" {
		logger.Debugf("no AES key provided, using default environment variable %s", DefaultVaultKeyEnv)
		keyEnv = DefaultVaultKeyEnv
	} else {
		logger.Debugf("using AES key from environment variable %s", keyEnv)
	}

	key := os.Getenv(keyEnv)
	if key == "" {
		key = generateAESKey(logger, keyEnv, logLevel)
		// this key needs to be set when initializing the vault
		if err := os.Setenv(keyEnv, key); err != nil {
			logger.FatalErr(fmt.Errorf("unable to set environment variable %s: %w", keyEnv, err))
		}
	} else {
		logger.Debugf("using existing AES key from environment variable %s", keyEnv)
	}

	storagePath = utils.ExpandPath(logger, storagePath, CacheDirectory(""), nil)
	if storagePath == "" {
		logger.Fatalf("unable to expand storage path: %s", storagePath)
	}

	opts := []vault.Option{
		vault.WithAESPath(storagePath),
		vault.WithProvider(vault.ProviderTypeAES256),
		vault.WithAESKeyFromEnv(keyEnv),
	}

	if keyFile != "" {
		keyFile = utils.ExpandPath(logger, keyFile, CacheDirectory(""), nil)
		if keyFile == "" {
			logger.Fatalf("unable to expand key file path: %s", keyFile)
		}
		opts = append(opts, vault.WithAESKeyFromFile(keyFile))
		if err := writeKeyToFile(logger, key, keyFile); err != nil {
			logger.Warnx("unable to write key to file", "err", err)
		}
	}

	v, cfg, err := vault.New(name, opts...)
	if err != nil {
		logger.FatalErr(err)
	}

	cfgPath := ConfigFilePath(v.ID())
	if err = vault.SaveConfigJSON(*cfg, cfgPath); err != nil {
		logger.FatalErr(fmt.Errorf("unable to save vault config: %w", err))
	}

	logger.PlainTextSuccess(fmt.Sprintf("Vault '%s' with AES256 encryption created successfully", v.ID()))
}

func generateAESKey(logger io.Logger, keyEnv, logLevel string) string {
	key, err := vault.GenerateEncryptionKey()
	if err != nil {
		logger.FatalErr(err)
	}

	if logLevel != "fatal" {
		logger.PlainTextSuccess(fmt.Sprintf("Your vault encryption key is: %s", key))
		newKeyMsg := fmt.Sprintf(
			"You will need this key to modify your vault data. Store it somewhere safe!\n"+
				"Set this value to the %s environment variable to access the vault in the future.\n",
			keyEnv,
		)
		logger.PlainTextInfo(newKeyMsg)
	} else {
		logger.PlainTextSuccess(fmt.Sprintf("Encryption key: %s", key))
	}
	return key
}

func NewAgeVault(logger io.Logger, name, storagePath, recipients, identityKey, identityFile string) {
	storagePath = utils.ExpandPath(logger, storagePath, CacheDirectory(""), nil)
	if storagePath == "" {
		logger.Fatalf("unable to expand storage path: %s", storagePath)
	}

	opts := []vault.Option{vault.WithAgePath(storagePath), vault.WithProvider(vault.ProviderTypeAge)}
	if recipients != "" {
		opts = append(opts, vault.WithAgeRecipients(strings.Split(recipients, ",")...))
	}
	if identityKey != "" {
		opts = append(opts, vault.WithAgeIdentityFromEnv(identityKey))
	}
	if identityFile != "" {
		identityFile = utils.ExpandPath(logger, identityFile, CacheDirectory(""), nil)
		opts = append(opts, vault.WithAgeIdentityFromFile(identityFile))
	}

	if identityKey == "" && identityFile == "" {
		logger.Debugf("no Age identity provided, using default environment variable %s", DefaultVaultIdentityEnv)
		opts = append(opts, vault.WithAgeIdentityFromEnv(DefaultVaultIdentityEnv))
	}

	v, cfg, err := vault.New(name, opts...)
	if err != nil {
		logger.FatalErr(err)
	}

	cfgPath := ConfigFilePath(v.ID())
	if err = vault.SaveConfigJSON(*cfg, cfgPath); err != nil {
		logger.FatalErr(fmt.Errorf("unable to save vault config: %w", err))
	}

	logger.PlainTextSuccess(fmt.Sprintf("Vault '%s' with Age encryption created successfully", v.ID()))
}

func VaultFromName(name string) (*VaultConfig, Vault, error) {
	if name == "" {
		return nil, nil, fmt.Errorf("vault name cannot be empty")
	} else if strings.ToLower(name) == DemoVaultReservedName {
		return newDemoVaultConfig(), newDemoVault(), nil
	}

	cfgPath := ConfigFilePath(name)
	cfg, err := vault.LoadConfigJSON(cfgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load vault config: %w", err)
	}

	switch cfg.Type {
	case vault.ProviderTypeAge:
		provider, err := vault.NewAgeVault(&cfg)
		return &cfg, provider, err
	case vault.ProviderTypeAES256:
		provider, err := vault.NewAES256Vault(&cfg)
		return &cfg, provider, err
	default:
		return nil, nil, fmt.Errorf("unsupported vault type: %s", cfg.Type)
	}
}

func CacheDirectory(subPath string) string {
	return filepath.Join(filesystem.CachedDataDirPath(), v2CacheDataDir, subPath)
}

func ConfigFilePath(vaultName string) string {
	return filepath.Join(
		filesystem.CachedDataDirPath(),
		v2CacheDataDir,
		fmt.Sprintf("configs/%s.json", vaultName),
	)
}

func writeKeyToFile(logger io.Logger, key, filePath string) error {
	if key == "" {
		return fmt.Errorf("no key provided to write to file")
	}
	if filePath == "" {
		return fmt.Errorf("no file path provided to write key")
	}

	if _, err := os.Stat(filePath); err == nil {
		logger.Debugf("key file already exists at %s, skipping write", filePath)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
		return fmt.Errorf("unable to create directory for key file: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(key), 0600); err != nil {
		return fmt.Errorf("unable to write key to file: %w", err)
	}
	logger.Infof("Key written to file: %s", filePath)

	return nil
}
