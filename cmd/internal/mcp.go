package internal

import (
	"github.com/spf13/cobra"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/services/mcp"
)

func RegisterMCPCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start the local model context protocol server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			mcpFunc(ctx, cmd, args)
		},
	}
	rootCmd.AddCommand(subCmd)
}

func mcpFunc(_ *context.Context, _ *cobra.Command, _ []string) {
	server := mcp.NewMCPServer()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
