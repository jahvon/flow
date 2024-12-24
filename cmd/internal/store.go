package internal

import (
	"fmt"
	"strings"

	"github.com/jahvon/tuikit/views"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/services/store"
)

func RegisterStoreCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:   "store",
		Short: "Manage the data store for persisting key-value data.",
		Long: "Manage the flow data store - a key-value store that persists data within and across executable runs. " +
			"Values set outside executables persist globally, while values set within executables persist only for " +
			"that execution scope.",
		Args: cobra.NoArgs,
	}
	registerStoreSetCmd(ctx, subCmd)
	registerStoreGetCmd(ctx, subCmd)
	registerStoreClearCmd(ctx, subCmd)
	rootCmd.AddCommand(subCmd)
}

func registerStoreSetCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:   "set KEY [VALUE]",
		Short: "Set a key-value pair in the store.",
		Long:  dataStoreDescription + "This will overwrite any existing value for the key.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			storeSetFunc(ctx, cmd, args)
		},
	}
	rootCmd.AddCommand(subCmd)
}

func storeSetFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	key := args[0]

	var value string
	switch {
	case len(args) == 1:
		form, err := views.NewForm(
			io.Theme(ctx.Config.Theme.String()),
			ctx.StdIn(),
			ctx.StdOut(),
			&views.FormField{
				Key:   "value",
				Type:  views.PromptTypeMultiline,
				Title: "Enter the value to store",
			})
		if err != nil {
			ctx.Logger.FatalErr(err)
		}
		if err = form.Run(ctx.Ctx); err != nil {
			ctx.Logger.FatalErr(err)
		}
		value = form.FindByKey("value").Value()
	case len(args) == 2:
		value = args[1]
	default:
		ctx.Logger.PlainTextWarn(fmt.Sprintf("merging multiple (%d) arguments into a single value", len(args)-1))
		value = strings.Join(args[1:], " ")
	}

	s, err := store.NewStore()
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	if err = s.CreateBucket(store.EnvironmentBucket()); err != nil {
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
	ctx.Logger.PlainTextInfo(fmt.Sprintf("Key %q set in the store", key))
}

func registerStoreGetCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "get KEY",
		Aliases: []string{"view"},
		Short:   "Get a value from the store by its key.",
		Long:    dataStoreDescription + "This will retrieve the value for the given key.",
		Args:    cobra.ExactArgs(1),
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
	if _, err = s.CreateAndSetBucket(store.EnvironmentBucket()); err != nil {
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
		Short:   "Clear data from the store. Use --full to remove all stored data.",
		Long:    dataStoreDescription + "This will remove all keys and values from the data store.",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			storeClearFunc(ctx, cmd, args)
		},
	}
	RegisterFlag(ctx, subCmd, *flags.StoreFullFlag)
	rootCmd.AddCommand(subCmd)
}

func storeClearFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	full := flags.ValueFor[bool](ctx, cmd, *flags.StoreFullFlag, false)
	if full {
		if err := store.DestroyStore(); err != nil {
			ctx.Logger.FatalErr(err)
		}
		ctx.Logger.PlainTextSuccess("Store store cleared")
		return
	}
	s, err := store.NewStore()
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			ctx.Logger.Error(err, "cleanup failure")
		}
	}()
	if err := s.DeleteBucket(store.EnvironmentBucket()); err != nil {
		ctx.Logger.FatalErr(err)
	}
	ctx.Logger.PlainTextSuccess("Store store cleared")
}

var dataStoreDescription = "The data store is a key-value store that can be used to persist data across executions. " +
	"Values that are set outside of an executable will persist across all executions until they are cleared. " +
	"When set within an executable, the data will only persist across serial or parallel sub-executables but all " +
	"values will be cleared when the parent executable completes.\n\n"
