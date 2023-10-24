//nolint:cyclop
package utils

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func ValidateAndGetContext(cmd *cobra.Command, set flags.FlagSet, currentConfig *config.RootConfig) (string, error) {
	var context string
	globalRaw, err := set.ValueFor(cmd, flags.ListGlobalContextFlag.Name)
	if err != nil {
		return "", fmt.Errorf("invalid global flag - %w", err)
	}
	global, _ := globalRaw.(bool)

	wsRaw, err := set.ValueFor(cmd, flags.FilterNamespaceFlag.Name)
	if err != nil {
		return "", fmt.Errorf("invalid namespace flag - %w", err)
	}
	ws, _ := wsRaw.(string)

	switch {
	case global && ws != "":
		return "", errors.New("cannot set a secret to both global and workspace scope")
	case global:
		context = "global"
	case ws != "":
		var found bool
		for _, curWs := range currentConfig.Workspaces {
			if curWs == ws {
				found = true
				break
			}
		}

		if !found {
			return "", fmt.Errorf("workspace %s does not exist", ws)
		}
		context = ws
	default:
		log.Debug().Msg("defaulting to the current workspace")
		context = currentConfig.CurrentWorkspace
	}

	return context, nil
}
