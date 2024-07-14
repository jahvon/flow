package utils

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd"
	"github.com/jahvon/flow/internal/context"
)

type CommandRunner struct {
	rootCmd *cobra.Command
}

func NewE2ECommandRunner(ctx *context.Context) *CommandRunner {
	rootCmd := cmd.NewRootCmd(ctx)
	return &CommandRunner{rootCmd: rootCmd}
}

func (r *CommandRunner) Run(ctx *context.Context, args ...string) error {
	r.rootCmd.SetArgs(args)
	if err := cmd.Execute(ctx, r.rootCmd); err != nil {
		return err
	}
	return nil
}
