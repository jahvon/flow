package internal

import (
	"fmt"
	"strconv"
	"strings"

	io2 "github.com/jahvon/tuikit/io"
	"github.com/jahvon/tuikit/views"
	extvault "github.com/jahvon/vault"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	vaultIO "github.com/jahvon/flow/internal/io/vault"
	"github.com/jahvon/flow/internal/vault"
	vaultV2 "github.com/jahvon/flow/internal/vault/v2"
	"github.com/jahvon/flow/types/config"
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
	registerEditVaultCmd(ctx, vaultCmd)
	registerMigrateVaultCmd(ctx, vaultCmd)
	// TODO: add command for testing vault connectivity
	rootCmd.AddCommand(vaultCmd)
}

func registerCreateVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	createCmd := &cobra.Command{
		Use:     "create NAME",
		Aliases: []string{"new", "add"},
		Short:   "Create a new vault.",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			vaultName := args[0]
			if vaultName == vaultV2.LegacyVaultReservedName || vaultName == vaultV2.DemoVaultReservedName {
				ctx.Logger.Fatalf("create is unsupported for the reserved vaults")
			} else if err := vault.ValidateReference(vaultName); err != nil {
				ctx.Logger.Fatalf("invalid vault name '%s': %v", vaultName, err)
			}

			if _, found := ctx.Config.Vaults[vaultName]; found {
				ctx.Logger.Fatalf("vault %s already exists", vaultName)
			}
		},
		Run: func(cmd *cobra.Command, args []string) { createVaultFunc(ctx, cmd, args) },
	}

	RegisterFlag(ctx, createCmd, *flags.VaultTypeFlag)
	RegisterFlag(ctx, createCmd, *flags.VaultPathFlag)
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
	vaultType := flags.ValueFor[string](ctx, cmd, *flags.VaultTypeFlag, false)
	vaultPath := flags.ValueFor[string](ctx, cmd, *flags.VaultPathFlag, false)

	switch strings.ToLower(vaultType) {
	case "aes256":
		keyEnv := flags.ValueFor[string](ctx, cmd, *flags.VaultKeyEnvFlag, false)
		keyFile := flags.ValueFor[string](ctx, cmd, *flags.VaultKeyFileFlag, false)
		logLevel := flags.ValueFor[string](ctx, cmd, *flags.LogLevel, false)
		vaultV2.NewAES256Vault(logger, vaultName, vaultPath, keyEnv, keyFile, logLevel)
	case "age":
		recipients := flags.ValueFor[string](ctx, cmd, *flags.VaultRecipientsFlag, false)
		identityEnv := flags.ValueFor[string](ctx, cmd, *flags.VaultIdentityEnvFlag, false)
		identityFile := flags.ValueFor[string](ctx, cmd, *flags.VaultIdentityFileFlag, false)
		vaultV2.NewAgeVault(logger, vaultName, vaultPath, recipients, identityEnv, identityFile)
	default:
		logger.Fatalf("unsupported vault type: %s - must be one of 'aes256' or 'age'", vaultType)
	}

	if ctx.Config.Vaults == nil {
		ctx.Config.Vaults = make(map[string]string)
	}
	ctx.Config.Vaults[vaultName] = vaultPath
	if err := filesystem.WriteConfig(ctx.Config); err != nil {
		logger.FatalErr(fmt.Errorf("unable to save user configuration: %w", err))
	}
}

func registerGetVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	getCmd := &cobra.Command{
		Use:     "get NAME",
		Aliases: []string{"view", "show"},
		Short:   "Get the details of a vault.",
		Args:    cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return vaultNames(ctx.Config), cobra.ShellCompDirectiveNoFileComp
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			var vaultName string
			if len(args) == 0 {
				vaultName = ctx.Config.CurrentVaultName()
			} else {
				vaultName = args[0]
			}

			if vaultName == vaultV2.LegacyVaultReservedName {
				ctx.Logger.Fatalf("get is unsupported for the legacy vault")
			} else if err := vaultV2.ValidateIdentifier(vaultName); err != nil {
				ctx.Logger.Fatalf("invalid vault name '%s': %v", vaultName, err)
			}

			StartTUI(ctx, cmd)
		},
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { getVaultFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, getCmd, *flags.OutputFormatFlag)
	vaultCmd.AddCommand(getCmd)
}

func getVaultFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)

	var vaultName string
	if len(args) == 0 {
		vaultName = ctx.Config.CurrentVaultName()
	} else {
		vaultName = args[0]
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
		PreRun: func(cmd *cobra.Command, args []string) {
			StartTUI(ctx, cmd)
		},
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { listVaultsFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, listCmd, *flags.OutputFormatFlag)
	vaultCmd.AddCommand(listCmd)
}

func listVaultsFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)

	cfg := ctx.Config
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
			return vaultNames(ctx.Config), cobra.ShellCompDirectiveNoFileComp
		},
		PreRun: func(cmd *cobra.Command, args []string) { validateVaults(ctx.Config, ctx.Logger) },
		Run:    func(cmd *cobra.Command, args []string) { removeVaultFunc(ctx, cmd, args) },
	}
	vaultCmd.AddCommand(removeCmd)
}

func removeVaultFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	vaultName := args[0]

	if vaultName == vaultV2.LegacyVaultReservedName || vaultName == vaultV2.DemoVaultReservedName {
		logger.Fatalf("remove is unsupported for the current vault")
	}

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
			return vaultNames(ctx.Config), cobra.ShellCompDirectiveNoFileComp
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			vaultName := args[0]
			reservedName := vaultName == vaultV2.LegacyVaultReservedName || vaultName == vaultV2.DemoVaultReservedName
			if reservedName {
				return
			}
			validateVaults(ctx.Config, ctx.Logger)
			if _, found := ctx.Config.Vaults[vaultName]; !found {
				ctx.Logger.Fatalf("vault %s not found", vaultName)
			}
		},
		Run: func(cmd *cobra.Command, args []string) { switchVaultFunc(ctx, cmd, args) },
	}
	vaultCmd.AddCommand(switchCmd)
}

func switchVaultFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	vaultName := args[0]
	userConfig := ctx.Config
	userConfig.CurrentVault = &vaultName

	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess("Vault set to " + vaultName)
}

func registerEditVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	editCmd := &cobra.Command{
		Use:     "edit NAME",
		Aliases: []string{"update", "modify"},
		Short:   "Edit the configuration of an existing vault.",
		Long: "Edit the configuration of an existing vault. " +
			"Note: You cannot change the vault type after creation.",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return vaultNames(ctx.Config), cobra.ShellCompDirectiveNoFileComp
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			validateVaults(ctx.Config, ctx.Logger)
			vaultName := args[0]
			if vaultName == vaultV2.LegacyVaultReservedName || vaultName == vaultV2.DemoVaultReservedName {
				ctx.Logger.Fatalf("edit is unsupported for the current vault")
			} else if err := vaultV2.ValidateIdentifier(vaultName); err != nil {
				ctx.Logger.Fatalf("invalid vault name '%s': %v", vaultName, err)
			}

			userConfig := ctx.Config
			if _, found := userConfig.Vaults[vaultName]; !found {
				ctx.Logger.Fatalf("vault %s not found", vaultName)
			}
		},
		Run: func(cmd *cobra.Command, args []string) { editVaultFunc(ctx, cmd, args) },
	}

	RegisterFlag(ctx, editCmd, *flags.VaultPathFlag)
	// AES flags
	RegisterFlag(ctx, editCmd, *flags.VaultKeyEnvFlag)
	RegisterFlag(ctx, editCmd, *flags.VaultKeyFileFlag)
	// Age flags
	RegisterFlag(ctx, editCmd, *flags.VaultRecipientsFlag)
	RegisterFlag(ctx, editCmd, *flags.VaultIdentityEnvFlag)
	RegisterFlag(ctx, editCmd, *flags.VaultIdentityFileFlag)

	vaultCmd.AddCommand(editCmd)
}

func editVaultFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	vaultName := args[0]

	vaultPath := flags.ValueFor[string](ctx, cmd, *flags.VaultPathFlag, false)
	keyEnv := flags.ValueFor[string](ctx, cmd, *flags.VaultKeyEnvFlag, false)
	keyFile := flags.ValueFor[string](ctx, cmd, *flags.VaultKeyFileFlag, false)
	recipients := flags.ValueFor[string](ctx, cmd, *flags.VaultRecipientsFlag, false)
	identityEnv := flags.ValueFor[string](ctx, cmd, *flags.VaultIdentityEnvFlag, false)
	identityFile := flags.ValueFor[string](ctx, cmd, *flags.VaultIdentityFileFlag, false)

	cfgPath := vaultV2.ConfigFilePath(vaultName)
	existingCfg, err := extvault.LoadConfigJSON(cfgPath)
	if err != nil {
		logger.Fatalf("failed to load vault configuration: %v", err)
	}

	// TODO: add support for appending KeySources and IdentitySources instead of overwriting them
	switch existingCfg.Type {
	case extvault.ProviderTypeAES256:
		if vaultPath != "" {
			existingCfg.Aes.StoragePath = vaultPath
		}
		if keyEnv != "" {
			existingCfg.Aes.KeySource = []extvault.KeySource{{
				Type: "env",
				Name: keyEnv,
			}}
		}
		if keyFile != "" {
			existingCfg.Aes.KeySource = []extvault.KeySource{{
				Type: "file",
				Path: keyFile,
			}}
		}
	case extvault.ProviderTypeAge:
		if vaultPath != "" {
			existingCfg.Age.StoragePath = vaultPath
		}
		if recipients != "" {
			existingCfg.Age.Recipients = strings.Split(recipients, ",")
		}
		if identityEnv != "" {
			existingCfg.Age.IdentitySources = []extvault.IdentitySource{{
				Type: "env",
				Name: identityEnv,
			}}
		}
		if identityFile != "" {
			existingCfg.Age.IdentitySources = []extvault.IdentitySource{{
				Type: "file",
				Path: identityFile,
			}}
		}
	default:
		logger.Fatalf("unsupported vault type: %s", existingCfg.Type)
	}

	if err = extvault.SaveConfigJSON(existingCfg, cfgPath); err != nil {
		logger.Fatalf("failed to save vault configuration: %v", err)
	}

	logger.PlainTextSuccess(fmt.Sprintf("Vault '%s' configuration updated successfully", vaultName))
}

func registerMigrateVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	migrateCmd := &cobra.Command{
		Use:   "migrate TARGET",
		Short: "Migrate the legacy vault to a newer vault.",
		Long: "Migrate the legacy vault to a newer vault type. " +
			"The target vault must exist and the encryption key must be set for the legacy vault. " +
			"Note: This will not remove the legacy vault, but will copy its contents to the target vault.",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return vaultNames(ctx.Config), cobra.ShellCompDirectiveNoFileComp
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			validateVaults(ctx.Config, ctx.Logger)
			vaultName := args[0]
			if vaultName == vaultV2.LegacyVaultReservedName || vaultName == vaultV2.DemoVaultReservedName {
				ctx.Logger.Fatalf("migrate is unsupported for the current vault")
			} else if err := vaultV2.ValidateIdentifier(vaultName); err != nil {
				ctx.Logger.Fatalf("invalid vault name '%s': %v", vaultName, err)
			}

			userConfig := ctx.Config
			if _, found := userConfig.Vaults[vaultName]; !found {
				ctx.Logger.Fatalf("vault %s not found", vaultName)
			}
		},
		Run: func(cmd *cobra.Command, args []string) { migrateVaultFunc(ctx, cmd, args) },
	}

	vaultCmd.AddCommand(migrateCmd)
}

func migrateVaultFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	targetVaultName := args[0]

	setAuthEnv(ctx, cmd, nil, true)
	legacyVault := vault.NewVault(logger)
	_, targetVault, err := vaultV2.VaultFromName(targetVaultName)
	if err != nil {
		logger.Fatalf("failed to load target vault '%s': %v", targetVaultName, err)
	}
	defer targetVault.Close()

	s1, err := legacyVault.GetAllSecrets()
	if err != nil {
		logger.Fatalf("failed to retrieve secrets from legacy vault: %v", err)
	}
	for name, secret := range s1 {
		if err := targetVault.SetSecret(name, vaultV2.NewSecretValue([]byte(secret.PlainTextString()))); err != nil {
			logger.Fatalf("failed to migrate secret '%s' to target vault '%s': %v", name, targetVaultName, err)
		}
	}

	logger.PlainTextSuccess(fmt.Sprintf("Legacy vault migrated to '%s'", targetVaultName))
}

func vaultNames(cfg *config.Config) []string {
	names := []string{vaultV2.LegacyVaultReservedName, vaultV2.DemoVaultReservedName}
	if cfg == nil || cfg.Vaults == nil {
		return nil
	}
	for name := range cfg.Vaults {
		names = append(names, name)
	}
	return names
}

func validateVaults(cfg *config.Config, logger io2.Logger) {
	if cfg == nil || cfg.Vaults == nil {
		logger.Fatalf("no vaults configured")
	}
}
