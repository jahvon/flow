package main

import (
	stdCtx "context"

	"github.com/flowexec/flow/cmd"
	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/io"
)

func main() {
	ctx := context.NewContext(stdCtx.Background(), io.Stdin, io.Stdout)
	defer ctx.Finalize()

	if ctx == nil {
		panic("failed to initialize context")
	}
	rootCmd := cmd.NewRootCmd(ctx)
	ctx.Ctx, ctx.CancelFunc = stdCtx.WithCancel(ctx.Ctx)
	if err := cmd.Execute(ctx, rootCmd); err != nil {
		ctx.Logger.FatalErr(err)
	}
}
