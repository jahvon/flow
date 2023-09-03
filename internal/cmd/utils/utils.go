package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/cmd/flags"
	"github.com/jahvon/tbox/internal/config"
)

func ValidateAndGetContext(cmd *cobra.Command, currentConfig *config.RootConfig) (string, error) {
	var err error
	var global *bool
	var ws *string
	var context string

	*global, err = strconv.ParseBool(cmd.Flag(flags.GlobalContextFlagName).Value.String())
	if err != nil {
		return "", fmt.Errorf("invalid global flag - %v", err)
	}
	*ws = cmd.Flag(flags.WorkspaceContextFlagName).Value.String()
	globalSet := global != nil
	wsSet := ws != nil

	if (globalSet && *global) && (wsSet && *ws != "") {
		return "", errors.New("cannot set a secret to both global and workspace scope")
	} else if globalSet && *global {
		context = "global"
	} else if wsSet && *ws != "" {
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
	} else {
		log.Debug().Msg("defaulting to the current workspace")
		context = currentConfig.CurrentWorkspace
	}

	return context, nil
}
