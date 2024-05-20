package internal

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
)

func RegisterFlag(ctx *context.Context, cmd *cobra.Command, flag flags.Metadata) {
	flagSet, err := flags.ToPflag(cmd, flag, false)
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	cmd.Flags().AddFlagSet(flagSet)
}

func RegisterPersistentFlag(ctx *context.Context, cmd *cobra.Command, flag flags.Metadata) {
	flagSet, err := flags.ToPflag(cmd, flag, true)
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	cmd.PersistentFlags().AddFlagSet(flagSet)
}

func MarkFlagRequired(ctx *context.Context, cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagRequired(name); err != nil {
		ctx.Logger.FatalErr(err)
	}
}

func MarkFlagMutuallyExclusive(cmd *cobra.Command, names ...string) {
	cmd.MarkFlagsMutuallyExclusive(names...)
}

func MarkOneFlagRequired(cmd *cobra.Command, names ...string) {
	cmd.MarkFlagsOneRequired(names...)
}

func MarkFlagFilename(ctx *context.Context, cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagFilename(name); err != nil {
		ctx.Logger.FatalErr(err)
	}
}
