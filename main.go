package main

import (
	stdCtx "context"
	"os"

	"github.com/jahvon/flow/cmd"
	"github.com/jahvon/flow/internal/context"
)

func main() {
	ctx := context.NewContext(stdCtx.Background(), os.Stdin, os.Stdout)
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
