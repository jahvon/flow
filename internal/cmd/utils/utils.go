//nolint:cyclop
package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func ValidateAndGetContext(cmd *cobra.Command, currentConfig *config.RootConfig) (string, error) {
	var global *bool
	var ws *string
	var context string

	globalFlag := cmd.Flag(flags.GlobalContextFlagName)
	wsFlag := cmd.Flag(flags.WorkspaceContextFlagName)
	if globalFlag != nil {
		val, err := strconv.ParseBool(cmd.Flag(flags.GlobalContextFlagName).Value.String())
		if err != nil {
			return "", fmt.Errorf("invalid global flag - %w", err)
		}
		global = &val
	}
	if wsFlag != nil {
		val := cmd.Flag(flags.WorkspaceContextFlagName).Value.String()
		ws = &val
	}
	globalSet := global != nil
	wsSet := ws != nil

	switch {
	case (globalSet && *global) && (wsSet && *ws != ""):
		return "", errors.New("cannot set a secret to both global and workspace scope")
	case globalSet && *global:
		context = "global"
	case wsSet && *ws != "":
		var found bool
		for _, curWs := range currentConfig.Workspaces {
			if curWs == *ws {
				found = true
				break
			}
		}

		if !found {
			return "", fmt.Errorf("workspace %s does not exist", *ws)
		}
		context = *ws
	default:
		log.Debug().Msg("defaulting to the current workspace")
		context = currentConfig.CurrentWorkspace
	}

	return context, nil
}
