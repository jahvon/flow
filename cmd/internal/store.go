package internal

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/store"
)

func RegisterStoreCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:   "store",
		Short: "Manage the data store.",
		Args:  cobra.NoArgs,
	}
	registerStoreSetCmd(ctx, subCmd)
	registerStoreGetCmd(ctx, subCmd)
	registerStoreClearCmd(ctx, subCmd)
	rootCmd.AddCommand(subCmd)
}

func registerStoreSetCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:   "set",
		Short: "Set a key-value pair in the data store.",
		Long:  dataStoreDescription + "This will overwrite any existing value for the key.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			storeSetFunc(ctx, cmd, args)
		},
	}
	rootCmd.AddCommand(subCmd)
}

func storeSetFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	s, err := store.NewStore()
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	if err = s.CreateBucket(); err != nil {
		ctx.Logger.FatalErr(err)
	}
	defer func() {
		if err = s.Close(); err != nil {
			ctx.Logger.Error(err, "cleanup failure")
		}
	}()
	if err = s.Set(key, value); err != nil {
		ctx.Logger.FatalErr(err)
	}
	ctx.Logger.Infof("key %q set in the flow store", key)
}

func registerStoreGetCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a value from the data store.",
		Long:  dataStoreDescription + "This will retrieve the value for the given key.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			storeGetFunc(ctx, cmd, args)
		},
	}
	rootCmd.AddCommand(subCmd)
}

func storeGetFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	key := args[0]

	s, err := store.NewStore()
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	if err = s.CreateBucket(); err != nil {
		ctx.Logger.FatalErr(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			ctx.Logger.Error(err, "cleanup failure")
		}
	}()
	value, err := s.Get(key)
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	ctx.Logger.PlainTextSuccess(value)
}

func registerStoreClearCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "clear",
		Aliases: []string{"reset"},
		Short:   "Clear the data store.",
		Long:    dataStoreDescription + "This will remove all keys and values from the data store.",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			storeClearFunc(ctx, cmd, args)
		},
	}
	rootCmd.AddCommand(subCmd)
}

func storeClearFunc(ctx *context.Context, _ *cobra.Command, _ []string) {
	s, err := store.NewStore()
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	if err := s.DeleteBucket(); err != nil {
		ctx.Logger.FatalErr(err)
	}
	ctx.Logger.PlainTextSuccess("data store cleared")
}

var dataStoreDescription = "The data store is a key-value store that can be used to persist data across executions. " +
	"Values that are set outside of an executable will persist across all executions until they are cleared. " +
	"When set within an executable, the data will only persist across serial or parallel sub-executables but all " +
	"values will be cleared when the parent executable completes.\n\n"
