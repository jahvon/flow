package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
)

var (
	log    = io.Log()
	curCtx *context.Context
)

func interactiveUIEnabled(cmd *cobra.Command) bool {
	disabled := getFlagValue[bool](cmd, *flags.NonInteractiveFlag)
	return !disabled && curCtx.UserConfig.InteractiveUI
}
