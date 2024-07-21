package utils

import (
	"github.com/jahvon/flow/cmd"
	"github.com/jahvon/flow/internal/context"
)

type CommandRunner struct{}

func NewE2ECommandRunner() *CommandRunner {
	return &CommandRunner{}
}

func (r *CommandRunner) Run(ctx *context.Context, args ...string) error {
	rootCmd := cmd.NewRootCmd(ctx)
	rootCmd.SetArgs(args)
	if err := cmd.Execute(ctx, rootCmd); err != nil {
		return err
	}
	return nil
}
