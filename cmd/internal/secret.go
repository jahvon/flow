package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jahvon/tuikit/views"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/crypto"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/secret"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterSecretCmd(ctx *context.Context, rootCmd *cobra.Command) {
	secretCmd := &cobra.Command{
		Use:     "secret",
		Aliases: []string{"scrt", "s"},
		Short:   "Manage flow secrets.",
	}
	registerInitVaultCmd(ctx, secretCmd)
	registerSetSecretCmd(ctx, secretCmd)
	registerListSecretCmd(ctx, secretCmd)
	registerGetSecretCmd(ctx, secretCmd)
	registerRemoveSecretCmd(ctx, secretCmd)
	rootCmd.AddCommand(secretCmd)
}

func registerRemoveSecretCmd(ctx *context.Context, secretCmd *cobra.Command) {
	removeCmd := &cobra.Command{
		Use:     "remove NAME",
		Aliases: []string{"rm", "delete", "destroy"},
		Short:   "Remove a secret from the vault.",
		Args:    cobra.ExactArgs(1),
		Run:     func(cmd *cobra.Command, args []string) { removeSecretFunc(ctx, cmd, args) },
	}
	secretCmd.AddCommand(removeCmd)
}

func removeSecretFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	reference := args[0]

	form, err := views.NewForm(
		io.Theme(),
		ctx.StdIn(),
		ctx.StdOut(),
		&views.FormField{
			Key:   "confirm",
			Type:  views.PromptTypeConfirm,
			Title: fmt.Sprintf("Are you sure you want to remove the secret '%s'?", reference),
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

	v := vault.NewVault(logger)
	if err = v.DeleteSecret(reference); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Secret %s removed from vault", reference))
}

func registerInitVaultCmd(ctx *context.Context, secretCmd *cobra.Command) {
	vaultCmd := &cobra.Command{
		Use:   "vault",
		Short: "Create a new flow secret vault.",
		Args:  cobra.NoArgs,
		Run:   func(cmd *cobra.Command, args []string) { initVaultFunc(ctx, cmd, args) },
	}
	secretCmd.AddCommand(vaultCmd)
}

func initVaultFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	generatedKey, err := crypto.GenerateKey()
	if err != nil {
		logger.FatalErr(err)
	}
	if err = vault.RegisterEncryptionKey(generatedKey); err != nil {
		logger.FatalErr(err)
	}

	if verbosity := flags.ValueFor[int](ctx, cmd, *flags.VerbosityFlag, false); verbosity >= 0 {
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

func registerSetSecretCmd(ctx *context.Context, secretCmd *cobra.Command) {
	setCmd := &cobra.Command{
		Use:     "set NAME [VALUE]",
		Aliases: []string{"new", "update"},
		Short:   "Update or create a secret in the flow secret vault.",
		Args:    cobra.MinimumNArgs(1),
		PreRun:  func(cmd *cobra.Command, args []string) { printContext(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { setSecretFunc(ctx, cmd, args) },
	}
	secretCmd.AddCommand(setCmd)
}

func setSecretFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	reference := args[0]

	var value string
	switch {
	case len(args) == 1:
		form, err := views.NewForm(
			io.Theme(),
			ctx.StdIn(),
			ctx.StdOut(),
			&views.FormField{
				Key:   "value",
				Type:  views.PromptTypeMasked,
				Title: "Enter the secret value",
			})
		if err != nil {
			logger.FatalErr(err)
		}
		if err := form.Run(ctx.Ctx); err != nil {
			logger.FatalErr(err)
		}
		value = form.FindByKey("value").Value()
	case len(args) == 2:
		value = args[1]
	default:
		logger.Warnx("merging multiple arguments into a single value", "count", len(args))
		value = strings.Join(args[1:], " ")
	}

	sv := vault.SecretValue(value)
	v := vault.NewVault(logger)
	err := v.SetSecret(reference, sv)
	if err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Secret %s set in vault", reference))
}

func registerListSecretCmd(ctx *context.Context, secretCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "View a list of secrets in the flow vault.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { listSecretFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, listCmd, *flags.OutputSecretAsPlainTextFlag)
	secretCmd.AddCommand(listCmd)
}

func listSecretFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	asPlainText := flags.ValueFor[bool](ctx, cmd, *flags.OutputSecretAsPlainTextFlag, false)

	v := vault.NewVault(logger)
	secrets, err := v.GetAllSecrets()
	if err != nil {
		logger.FatalErr(err)
	}

	interactiveUI := TUIEnabled(ctx, cmd)
	if interactiveUI {
		secret.LoadSecretListView(ctx, asPlainText)
	} else {
		for ref, s := range secrets {
			if asPlainText {
				logger.PlainTextInfo(fmt.Sprintf("%s: %s", ref, s.PlainTextString()))
			} else {
				logger.PlainTextInfo(fmt.Sprintf("%s: %s", ref, s.String()))
			}
		}
	}
}

func registerGetSecretCmd(ctx *context.Context, secretCmd *cobra.Command) {
	getCmd := &cobra.Command{
		Use:     "get REFERENCE",
		Aliases: []string{"find"},
		Short:   "Get the value of a secret in the secret vault.",
		Args:    cobra.ExactArgs(1),
		Run:     func(cmd *cobra.Command, args []string) { getSecretFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, getCmd, *flags.OutputSecretAsPlainTextFlag)
	RegisterFlag(ctx, getCmd, *flags.CopyFlag)
	secretCmd.AddCommand(getCmd)
}

func getSecretFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	reference := args[0]
	asPlainText := flags.ValueFor[bool](ctx, cmd, *flags.OutputSecretAsPlainTextFlag, false)
	copyValue := flags.ValueFor[bool](ctx, cmd, *flags.CopyFlag, false)

	v := vault.NewVault(logger)
	s, err := v.GetSecret(reference)
	if err != nil {
		logger.FatalErr(err)
	}

	if asPlainText {
		logger.PlainTextInfo(s.PlainTextString())
	} else {
		logger.PlainTextInfo(s.String())
	}

	if copyValue {
		if err := clipboard.WriteAll(s.PlainTextString()); err != nil {
			logger.Error(err, "\nunable to copy secret value to clipboard")
		} else {
			logger.PlainTextSuccess("\ncopied secret value to clipboard")
		}
	}
}
