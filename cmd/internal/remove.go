package internal

import (
	"fmt"
	"strconv"

	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterRemoveCmd(ctx *context.Context, rootCmd *cobra.Command) {
	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove a flow entity.",
	}
	registerRemoveWsCmd(ctx, removeCmd)
	registerRemoveSecretCmd(ctx, removeCmd)
	rootCmd.AddCommand(removeCmd)
}

func registerRemoveWsCmd(ctx *context.Context, removeCmd *cobra.Command) {
	wsCmd := &cobra.Command{
		Use:     "workspace NAME",
		Aliases: []string{"ws"},
		Short:   "Remove an existing workspace from the list of known workspaces.",
		Long: "Remove an existing workspace. File contents will remain in the corresponding directory but the " +
			"workspace will be unlinked from the flow global configurations.\nNote: You cannot remove the current workspace.",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return maps.Keys(ctx.UserConfig.Workspaces), cobra.ShellCompDirectiveNoFileComp
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			interactive.InitInteractiveCommand(ctx, cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			removeWsFunc(ctx, cmd, args)
		},
	}
	removeCmd.AddCommand(wsCmd)
}

func removeWsFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	name := args[0]

	inputs, err := components.ProcessInputs(io.Theme(), &components.TextInput{
		Key:    "confirm",
		Prompt: fmt.Sprintf("Are you sure you want to remove the workspace '%s'? (y/n)", name),
	})
	if err != nil {
		logger.FatalErr(err)
	}
	resp := inputs.FindByKey("confirm").Value()
	confirmed, _ := strconv.ParseBool(resp)
	if !confirmed {
		logger.Warnf("Aborting")
		return
	}

	userConfig := ctx.UserConfig
	if name == userConfig.CurrentWorkspace {
		logger.Fatalf("cannot remove the current workspace")
	}
	if _, found := userConfig.Workspaces[name]; !found {
		logger.Fatalf("workspace %s was not found", name)
	}

	delete(userConfig.Workspaces, name)
	if err := filesystem.WriteUserConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}

	logger.Warnf("Workspace '%s' removed", name)

	if err := cache.UpdateAll(logger); err != nil {
		logger.FatalErr(errors.Wrap(err, "unable to update cache"))
	}
}

func registerRemoveSecretCmd(ctx *context.Context, removeCmd *cobra.Command) {
	secretCmd := &cobra.Command{
		Use:     "secret <name>",
		Aliases: []string{"scrt"},
		Short:   "Remove a secret from the vault.",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			interactive.InitInteractiveCommand(ctx, cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			removeSecretFunc(ctx, cmd, args)
		},
	}
	removeCmd.AddCommand(secretCmd)
}

func removeSecretFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	reference := args[0]

	inputs, err := components.ProcessInputs(io.Theme(), &components.TextInput{
		Key:    "confirm",
		Prompt: fmt.Sprintf("Are you sure you want to remove the secret '%s'? (y/n)", reference),
	})
	if err != nil {
		logger.FatalErr(err)
	}
	resp := inputs.FindByKey("confirm").Value()
	confirmed, _ := strconv.ParseBool(resp)
	if !confirmed {
		logger.Warnf("Aborting")
		return
	}

	v := vault.NewVault(logger)
	if err = v.DeleteSecret(reference); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Secret %s removed from vault", reference))
}
