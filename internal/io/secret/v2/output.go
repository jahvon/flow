package secret

import (
	"fmt"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/io/common"
	vaultV2 "github.com/flowexec/flow/internal/vault/v2"
)

func PrintSecrets(ctx *context.Context, vaultName string, vlt vaultV2.Vault, format string, plaintext bool) {
	logger := ctx.Logger

	secrets, err := vaultV2.NewSecretList(vaultName, vlt)
	if err != nil {
		logger.FatalErr(err)
	}

	if plaintext {
		secrets = secrets.AsPlaintext()
	} else {
		secrets = secrets.AsObfuscatedText()
	}

	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := secrets.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal secrets - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), str)
	case common.JSONFormat:
		str, err := secrets.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal secrets - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), str)
	}
}
