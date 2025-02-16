package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	envPrefix           = "FLOW_VAULT_"
	vaultConfigFileName = "vaults"
	vaultConfigFileExt  = ".yaml"
)

type VaultConfig struct {
	// Provider specifies the provider's type
	Provider string `koanf:"provider"`
	// Config specifies the provider-specific configuration
	Config map[string]interface{} `koanf:"config"`
}

type Config struct {
	// Current specifies the name of thr current vault to use
	Current string `koanf:"current"`
	// Vaults specifies a map of vault names to their configurations
	Vaults map[string]VaultConfig `koanf:"vaults"`
}

type LoadOptions struct {
	ConfigPath        string
	AutoDiscoveryPath string
	AllowEnv          bool
	RequireConfig     bool
}

// LoadCOnfig loads the vault configuration from the specified sources
func LoadConfig(opts LoadOptions) (*Config, error) {
	k := koanf.New(".")

	if opts.AutoDiscoveryPath != "" {
		if err := loadFromDir(k, opts.AutoDiscoveryPath); err != nil {
			if opts.RequireConfig || !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("failed to load config dir: %w", err)
			}
		}
	}

	if opts.ConfigPath != "" {
		if err := loadFromFile(k, opts.ConfigPath); err != nil {
			if opts.RequireConfig || !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("failed to load config file: %w", err)
			}
		}
	}

	if opts.AllowEnv {
		if err := loadFromEnv(k); err != nil {
			return nil, fmt.Errorf("failed to load environment variables: %w", err)
		}
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// loadFromFile loads configuration from a file, supporting multiple formats
func loadFromFile(k *koanf.Koanf, path string) error {
	var parser koanf.Parser

	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		parser = yaml.Parser()
	case ".json":
		parser = json.Parser()
	case ".toml":
		parser = toml.Parser()
	default:
		return fmt.Errorf("unsupported config file format: %s", path)
	}

	if err := k.Load(file.Provider(path), parser); err != nil {
		return err
	}

	return nil
}

// loadFromDir loads configuration from a directory, supporting multiple formats
func loadFromDir(k *koanf.Koanf, path string) error {
	if _, err := os.Stat(filepath.Join(path, vaultConfigFileName+vaultConfigFileExt)); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return loadFromFile(k, filepath.Join(path, vaultConfigFileName+vaultConfigFileExt))
	}
	if _, err := os.Stat(filepath.Join(path, vaultConfigFileName+".yml")); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return loadFromFile(k, filepath.Join(path, vaultConfigFileName+".yml"))
	}
	if _, err := os.Stat(filepath.Join(path, vaultConfigFileName+".json")); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return loadFromFile(k, filepath.Join(path, vaultConfigFileName+".json"))
	}
	if _, err := os.Stat(filepath.Join(path, vaultConfigFileName+".toml")); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return loadFromFile(k, filepath.Join(path, vaultConfigFileName+".toml"))
	}
	return fmt.Errorf("no config file found with name %s: %w", vaultConfigFileName, os.ErrNotExist)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(k *koanf.Koanf) error {
	return k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
}
