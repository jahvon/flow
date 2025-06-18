package secret

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io/common"
	"github.com/jahvon/flow/internal/vault"
)

type secretOutput struct {
	Secrets map[string]string `json:"secrets" yaml:"secrets"`
}

func PrintSecrets(ctx *context.Context, secrets map[string]vault.SecretValue, format string, plaintext bool) {
	if secrets == nil {
		return
	}
	logger := ctx.Logger
	output := secretOutput{}
	for key, value := range secrets {
		if plaintext {
			output.Secrets[key] = value.PlainTextString()
		} else {
			output.Secrets[key] = value.ObfuscatedString()
		}
	}
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := yaml.Marshal(output)
		if err != nil {
			logger.Fatalf("Failed to marshal secrets - %v", err)
		}
		_, _ = fmt.Fprintf(ctx.StdOut(), string(str))
	case common.JSONFormat:
		str, err := yaml.Marshal(output)
		if err != nil {
			logger.Fatalf("Failed to marshal secrets - %v", err)
		}
		_, _ = fmt.Fprintf(ctx.StdOut(), string(str))
	}
}
