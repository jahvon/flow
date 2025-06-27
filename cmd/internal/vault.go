package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jahvon/tuikit/views"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	vaultIO "github.com/jahvon/flow/internal/io/vault"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterVaultCmd(ctx *context.Context, rootCmd *cobra.Command) {
	vaultCmd := &cobra.Command{
		Use:     "vault",
		Aliases: []string{"vlt", "vaults"},
		Short:   "Manage sensitive secret stores.",
		Args:    cobra.NoArgs,
	}
	registerCreateVaultCmd(ctx, vaultCmd)
	registerGetVaultCmd(ctx, vaultCmd)
	registerListVaultCmd(ctx, vaultCmd)
	registerSwitchVaultCmd(ctx, vaultCmd)
	registerRemoveVaultCmd(ctx, vaultCmd)
	rootCmd.AddCommand(vaultCmd)
}

func registerCreateVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	createCmd := &cobra.Command{
		Use:     "create NAME",
		Aliases: []string{"new", "add"},
		Short:   "Create a new vault.",
		Args:    cobra.ExactArgs(1),
		Run:     func(cmd *cobra.Command, args []string) { createVaultFunc(ctx, cmd, args) },
	}

	RegisterFlag(ctx, createCmd, *flags.VaultTypeFlag)
	RegisterFlag(ctx, createCmd, *flags.VaultPathFlag)
	RegisterFlag(ctx, createCmd, *flags.VaultInteractiveFlag)
	// AES flags
	RegisterFlag(ctx, createCmd, *flags.VaultKeyEnvFlag)
	RegisterFlag(ctx, createCmd, *flags.VaultKeyFileFlag)
	// Age flags
	RegisterFlag(ctx, createCmd, *flags.VaultRecipientsFlag)
	RegisterFlag(ctx, createCmd, *flags.VaultIdentityEnvFlag)
	RegisterFlag(ctx, createCmd, *flags.VaultIdentityFileFlag)

	vaultCmd.AddCommand(createCmd)
}

func createVaultFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger

	vaultName := args[0]
	if err := vault.ValidateReference(vaultName); err != nil {
		logger.Fatalf("invalid vault name '%s': %v", vaultName, err)
	}
	vaultType := flags.ValueFor[string](ctx, cmd, *flags.VaultTypeFlag, false)
	vaultPath := flags.ValueFor[string](ctx, cmd, *flags.VaultPathFlag, false)

	switch strings.ToLower(vaultType) {
	case "aes256":
		keyEnv := flags.ValueFor[string](ctx, cmd, *flags.VaultKeyEnvFlag, false)
		keyFile := flags.ValueFor[string](ctx, cmd, *flags.VaultKeyFileFlag, false)
		logLevel := flags.ValueFor[string](ctx, cmd, *flags.LogLevel, false)
		vault.NewAES256Vault(logger, vaultName, vaultPath, keyEnv, keyFile, logLevel)
	case "age":
		recipients := flags.ValueFor[string](ctx, cmd, *flags.VaultRecipientsFlag, false)
		identityEnv := flags.ValueFor[string](ctx, cmd, *flags.VaultIdentityEnvFlag, false)
		identityFile := flags.ValueFor[string](ctx, cmd, *flags.VaultIdentityFileFlag, false)
		vault.NewAgeVault(logger, vaultName, vaultPath, recipients, identityEnv, identityFile)
	default:
		logger.Fatalf("unsupported vault type: %s - must be one of 'aes256' or 'age'", vaultType)
	}
}

func registerGetVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	getCmd := &cobra.Command{
		Use:     "get NAME",
		Aliases: []string{"view", "show"},
		Short:   "Get the details of a vault.",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return maps.Keys(ctx.Config.Vaults), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) { getVaultFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, getCmd, *flags.OutputFormatFlag)
	vaultCmd.AddCommand(getCmd)
}

func getVaultFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)

	vaultName := args[0]
	if err := vault.ValidateReference(vaultName); err != nil {
		logger.Fatalf("invalid vault name '%s': %v", vaultName, err)
	}

	if TUIEnabled(ctx, cmd) {
		view := vaultIO.NewVaultView(ctx.TUIContainer, vaultName)
		SetView(ctx, cmd, view)
	} else {
		vaultIO.PrintVault(logger, outputFormat, vaultName)
	}
}

func registerListVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available vaults.",
		Args:    cobra.NoArgs,
		Run:     func(cmd *cobra.Command, args []string) { listVaultsFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, listCmd, *flags.OutputFormatFlag)
	vaultCmd.AddCommand(listCmd)
}

func listVaultsFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)

	cfg := ctx.Config
	if cfg.Vaults == nil || len(cfg.Vaults) == 0 {
		logger.Fatalf("no vaults configured")
	}

	if TUIEnabled(ctx, cmd) {
		view := vaultIO.NewVaultListView(ctx.TUIContainer, maps.Keys(cfg.Vaults))
		SetView(ctx, cmd, view)
	} else {
		vaultIO.PrintVaultList(logger, outputFormat, maps.Keys(cfg.Vaults))
	}
}

func registerRemoveVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	removeCmd := &cobra.Command{
		Use:     "remove NAME",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove an existing vault.",
		Long: "Remove an existing vault by its name. The vault data will remain in it's original location, " +
			"but the vault will be unlinked from the global configuration.\nNote: You cannot remove the current vault.",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return maps.Keys(ctx.Config.Vaults), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) { removeVaultFunc(ctx, cmd, args) },
	}
	vaultCmd.AddCommand(removeCmd)
}

func removeVaultFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	vaultName := args[0]

	form, err := views.NewForm(
		io.Theme(ctx.Config.Theme.String()),
		ctx.StdIn(),
		ctx.StdOut(),
		&views.FormField{
			Key:   "confirm",
			Type:  views.PromptTypeConfirm,
			Title: fmt.Sprintf("Are you sure you want to remove the vault '%s'?", vaultName),
		})
	if err != nil {
		logger.FatalErr(err)
	}
	if err := form.Run(ctx.Ctx); err != nil {
		logger.FatalErr(err)
	}
	resp := form.FindByKey("confirm").Value()
	if truthy, _ := strconv.ParseBool(resp); !truthy {
		logger.Warnf("Aborting")
		return
	}

	userConfig := ctx.Config
	if userConfig.CurrentVault != nil && vaultName == *userConfig.CurrentVault {
		logger.Fatalf("cannot remove the current vault")
	}
	if _, found := userConfig.Vaults[vaultName]; !found {
		logger.Fatalf("vault %s was not found", vaultName)
	}

	delete(userConfig.Vaults, vaultName)
	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}

	logger.Warnf("Vault '%s' deleted", vaultName)

}

func registerSwitchVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	switchCmd := &cobra.Command{
		Use:     "switch NAME",
		Aliases: []string{"use", "set"},
		Short:   "Switch the active vault.",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return maps.Keys(ctx.Config.Vaults), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) { switchVaultFunc(ctx, cmd, args) },
	}
	vaultCmd.AddCommand(switchCmd)
}

func switchVaultFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	vaultName := args[0]
	userConfig := ctx.Config
	if _, found := userConfig.Vaults[vaultName]; !found {
		logger.Fatalf("vault %s not found", vaultName)
	}
	userConfig.CurrentVault = &vaultName

	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess("Vault set to " + vaultName)
}
