package internal

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/crypto"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterVaultCmd(ctx *context.Context, secretCmd *cobra.Command) {
	vaultCmd := &cobra.Command{
		Use:     "vault",
		Aliases: []string{"vlt", "vaults"},
		Short:   "Manage sensitive secret stores.",
		Args:    cobra.NoArgs,
	}
	registerCreateVaultCmd(ctx, vaultCmd)
	registerListVaultCmd(ctx, vaultCmd)
	registerSwitchVaultCmd(ctx, vaultCmd)
	secretCmd.AddCommand(vaultCmd)
}

func registerCreateVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "add"},
		Short:   "Create a new vault.",
		Args:    cobra.NoArgs,
		Run:     func(cmd *cobra.Command, args []string) { createSecretVaultFunc(ctx, cmd, args) },
	}
	vaultCmd.AddCommand(createCmd)
}

func createSecretVaultFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	generatedKey, err := crypto.GenerateKey()
	if err != nil {
		logger.FatalErr(err)
	}
	if err = vault.RegisterEncryptionKey(generatedKey); err != nil {
		logger.FatalErr(err)
	}

	if verbosity := flags.ValueFor[int](ctx, cmd, *flags.LogLevel, false); verbosity >= 0 {
		logger.PlainTextSuccess(fmt.Sprintf("Your vault encryption key is: %s", generatedKey))
		newKeyMsg := fmt.Sprintf(
			"You will need this key to modify your vault data. Store it somewhere safe!\n"+
				"Set this value to the %s environment variable if you do not want to be prompted for it every time.",
			vault.EncryptionKeyEnvVar,
		)
		logger.PlainTextInfo(newKeyMsg)
	} else {
		logger.PlainTextSuccess(fmt.Sprintf("Encryption key: %s", generatedKey))
	}
}

func registerListVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available vaults.",
		Args:    cobra.NoArgs,
		Run:     func(cmd *cobra.Command, args []string) { ctx.Logger.Fatalf("not implemented yet") },
	}
	vaultCmd.AddCommand(listCmd)
}

func registerSwitchVaultCmd(ctx *context.Context, vaultCmd *cobra.Command) {
	switchCmd := &cobra.Command{
		Use:     "switch NAME",
		Aliases: []string{"use", "set"},
		Short:   "Switch the active vault.",
		Args:    cobra.ExactArgs(1),
		Run:     func(cmd *cobra.Command, args []string) { ctx.Logger.Fatalf("not implemented yet") },
	}
	vaultCmd.AddCommand(switchCmd)
}
