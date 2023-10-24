package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/crypto"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

var vaultCmd = &cobra.Command{
	Use:     "vault",
	Aliases: []string{"v"},
	GroupID: DataGroup.ID,
	Short:   "Manage Flow's secret vault data.",
}

var vaultCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a new vault.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		generatedKey, err := crypto.GenerateKey()
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		if err = vault.RegisterEncryptionKey(generatedKey); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintSuccess(fmt.Sprintf("Your vault encryption key is: %s", generatedKey))
		newKeyMsg := fmt.Sprintf(
			"You will need this key to modify your vault data. Store it somewhere safe!\n"+
				"Set this value to the %s environment variable if you do not want to be prompted for it every time.",
			vault.EncryptionKeyEnvVar,
		)
		io.PrintNotice(newKeyMsg)
	},
}

var vaultSetCmd = &cobra.Command{
	Use:     "set <name> <value>",
	Aliases: []string{"s"},
	Short:   "Set a secret in the vault.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		reference := args[0]
		value := args[1]

		secret := vault.Secret(value)
		v := vault.NewVault()
		err := v.SetSecret(reference, secret)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess(fmt.Sprintf("Secret %s set in vault", reference))
	},
}

var vaultGetCmd = &cobra.Command{
	Use:     "get <name>",
	Aliases: []string{"g"},
	Short:   "Get a secret from the vault.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reference := args[0]
		plainTextFlag, err := Flags.ValueFor(cmd, flags.OutputSecretAsPlainTextFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		asPlainText, _ := plainTextFlag.(bool)

		v := vault.NewVault()
		secret, err := v.GetSecret(reference)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		if asPlainText {
			io.PrintNotice(secret.PlainTextString())
		} else {
			io.PrintNotice(secret.String())
		}
	},
}

var vaultListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all secrets in the vault.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		plainTextFlag, err := Flags.ValueFor(cmd, flags.OutputSecretAsPlainTextFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		asPlainText, _ := plainTextFlag.(bool)

		v := vault.NewVault()
		secrets, err := v.GetAllSecrets()
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		for ref, secret := range secrets {
			if asPlainText {
				io.PrintNotice(fmt.Sprintf("%s: %s", ref, secret.PlainTextString()))
			} else {
				io.PrintNotice(fmt.Sprintf("%s: %s", ref, secret.String()))
			}
		}
	},
}

var vaultDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a secret from the vault.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reference := args[0]

		v := vault.NewVault()
		err := v.DeleteSecret(reference)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess(fmt.Sprintf("Secret %s deleted from vault", reference))
	},
}

func init() {
	vaultCmd.AddCommand(vaultCreateCmd)
	vaultCmd.AddCommand(vaultSetCmd)

	registerFlagOrPanic(vaultGetCmd, *flags.OutputSecretAsPlainTextFlag)
	vaultCmd.AddCommand(vaultGetCmd)

	registerFlagOrPanic(vaultListCmd, *flags.OutputSecretAsPlainTextFlag)
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultDeleteCmd)

	rootCmd.AddCommand(vaultCmd)
}
