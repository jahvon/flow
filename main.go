package main

import (
	stdCtx "context"

	"github.com/jahvon/flow/cmd"
	"github.com/jahvon/flow/internal/context"
)

func main() {
	ctx := context.NewContext(stdCtx.Background())
	defer ctx.Finalize()

	if ctx == nil {
		panic("failed to initialize context")
	}
	ctx.Ctx, ctx.CancelFunc = stdCtx.WithCancel(ctx.Ctx)
	cmd.Execute(ctx)
}
