package utils

import (
	"fmt"

	"github.com/jahvon/flow/cmd"
	"github.com/jahvon/flow/internal/context"
)

type CommandRunner struct{}

func NewE2ECommandRunner() *CommandRunner {
	return &CommandRunner{}
}

func (r *CommandRunner) Run(ctx *context.Context, args ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	rootCmd := cmd.NewRootCmd(ctx)
	rootCmd.SetArgs(args)
	rootCmd.SetIn(ctx.StdIn())
	rootCmd.SetOut(ctx.StdOut())
	if err = cmd.Execute(ctx, rootCmd); err != nil {
		return err
	}
	return
}
