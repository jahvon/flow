package vault

import (
	tuikitIO "github.com/flowexec/tuikit/io"

	"github.com/flowexec/flow/internal/io/common"
)

func PrintVault(logger tuikitIO.Logger, format, vaultName string) {
	if vaultName == "" {
		logger.Fatalf("Vault name was not provided")
	}

	vault, err := vaultFromName(vaultName)
	if err != nil {
		logger.Fatalf("Failed to get vault %s - %v", vaultName, err)
	}

	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := vault.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal user config - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := vault.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal user config - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintVaultList(logger tuikitIO.Logger, format string, vaultNames []string) {
	vaults := vaultCollection{Vaults: make([]*vaultEntity, 0, len(vaultNames))}
	for _, name := range vaultNames {
		vault, err := vaultFromName(name)
		if err != nil {
			logger.Fatalf("Vault error %s - %v", name, err)
		}
		vaults.Vaults = append(vaults.Vaults, vault)
	}
	logger.Debugf("listing %d vaults", len(vaults.Vaults))

	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := vaults.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal vault list - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := vaults.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal vault list - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
