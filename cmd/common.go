package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

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

func GenerateMarkdownTree(dir string) error {
	return doc.GenMarkdownTree(rootCmd, dir)
}
