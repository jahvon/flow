package secret

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/io/common"
	"github.com/flowexec/flow/internal/logger"
	"github.com/flowexec/flow/internal/vault"
)

type secretOutput struct {
	Secrets map[string]string `json:"secrets" yaml:"secrets"`
}

func PrintSecrets(ctx *context.Context, secrets map[string]vault.SecretValue, format string, plaintext bool) {
	if secrets == nil {
		return
	}
	output := secretOutput{
		Secrets: make(map[string]string, len(secrets)),
	}
	for key, value := range secrets {
		if plaintext {
			output.Secrets[key] = value.PlainTextString()
		} else {
			output.Secrets[key] = value.ObfuscatedString()
		}
	}
	// TODO: switch to using the SecretList type or something similar
	switch common.NormalizeFormat(format) {
	case common.YAMLFormat:
		str, err := yaml.Marshal(output)
		if err != nil {
			logger.Log().Fatalf("Failed to marshal secrets - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), string(str))
	case common.JSONFormat:
		str, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			logger.Log().Fatalf("Failed to marshal secrets - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), string(str))
	}
}
