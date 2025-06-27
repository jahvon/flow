package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jahvon/tuikit/io"
	"github.com/jahvon/vault"

	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/utils"
)

const (
	DefaultVaultKeyEnv = "FLOW_VAULT_KEY"

	v2CacheDataDir = "vaults"
)

func NewAES256Vault(logger io.Logger, name, storagePath, keyEnv, keyFile, logLevel string) {
	key, err := vault.GenerateEncryptionKey()
	if err != nil {
		logger.FatalErr(err)
	}

	if logLevel != "fatal" {
		logger.PlainTextSuccess(fmt.Sprintf("Your vault encryption key is: %s", key))
		newKeyMsg := fmt.Sprintf(
			"You will need this key to modify your vault data. Store it somewhere safe!\n"+
				"Set this value to the %s environment variable if you do not want to be prompted for it every time.",
			DefaultVaultKeyEnv,
		)
		logger.PlainTextInfo(newKeyMsg)
	} else {
		logger.PlainTextSuccess(fmt.Sprintf("Encryption key: %s", key))
	}

	if storagePath == "" {
		storagePath = CacheDirectory("")
	}

	opts := []vault.Option{vault.WithLocalPath(storagePath), vault.WithProvider(vault.ProviderTypeAES256)}
	if keyEnv != "" {
		opts = append(opts, vault.WithAESKeyFromEnv(keyEnv))
	}
	if keyFile != "" {
		opts = append(opts, vault.WithAESKeyFromFile(keyFile))
		if err := writeKeyToFile(logger, key, keyFile); err != nil {
			logger.Warnx("unable to write key to file", "err", err)
		}
	}

	if keyEnv == "" && keyFile == "" {
		logger.Debugf("no AES key provided, using default environment variable %s", DefaultVaultKeyEnv)
		opts = append(opts, vault.WithAESKeyFromEnv(DefaultVaultKeyEnv))
	}

	v, cfg, err := vault.New(name, opts...)
	if err != nil {
		logger.FatalErr(err)
	}

	cfgPath := CacheDirectory(fmt.Sprintf("configs/%s.json", v.ID()))
	if err = vault.SaveConfigJSON(*cfg, cfgPath); err != nil {
		logger.FatalErr(fmt.Errorf("unable to save vault config: %w", err))
	}

	logger.PlainTextSuccess(fmt.Sprintf("Vault '%s' with AES256 encryption created successfully", v.ID()))
}

func NewAgeVault(logger io.Logger, name, storagePath, recipients, identityKey, identityFile string) {
	if storagePath == "" {
		storagePath = CacheDirectory("")
	}

	opts := []vault.Option{vault.WithLocalPath(storagePath), vault.WithProvider(vault.ProviderTypeAge)}
	if recipients != "" {
		opts = append(opts, vault.WithAgeRecipients(strings.Split(recipients, ",")...))
	}
	if identityKey != "" {
		opts = append(opts, vault.WithAgeIdentityFromEnv(identityKey))
	}
	if identityFile != "" {
		opts = append(opts, vault.WithAgeIdentityFromFile(identityFile))
	}

	if identityKey == "" && identityFile == "" {
		logger.Debugf("no Age identity provided, using default environment variable %s", DefaultVaultKeyEnv)
		opts = append(opts, vault.WithAgeIdentityFromEnv(DefaultVaultKeyEnv))
	}

	v, cfg, err := vault.New(name, opts...)
	if err != nil {
		logger.FatalErr(err)
	}

	cfgPath := CacheDirectory(fmt.Sprintf("configs/%s.json", v.ID()))
	if err = vault.SaveConfigJSON(*cfg, cfgPath); err != nil {
		logger.FatalErr(fmt.Errorf("unable to save vault config: %w", err))
	}

	logger.PlainTextSuccess(fmt.Sprintf("Vault '%s' with Age encryption created successfully", v.ID()))
}

func CacheDirectory(subPath string) string {
	return filepath.Join(filesystem.CachedDataDirPath(), v2CacheDataDir, subPath)
}

func writeKeyToFile(logger io.Logger, key, filePath string) error {
	if key == "" {
		return fmt.Errorf("no key provided to write to file")
	}
	if filePath == "" {
		return fmt.Errorf("no file path provided to write key")
	}

	expandedPath := utils.ExpandPath(logger, filePath, "", nil)
	if expandedPath == "" {
		return fmt.Errorf("failed to expand path: %s", filePath)
	}
	if _, err := os.Stat(expandedPath); err == nil {
		logger.Debugf("key file already exists at %s, skipping write", expandedPath)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(expandedPath), 0750); err != nil {
		return fmt.Errorf("unable to create directory for key file: %w", err)
	}

	if err := os.WriteFile(expandedPath, []byte(key), 0600); err != nil {
		return fmt.Errorf("unable to write key to file: %w", err)
	}

	return nil
}
