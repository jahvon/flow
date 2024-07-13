package secret

import (
	"fmt"

	"github.com/jahvon/tuikit/components"
	"github.com/jahvon/tuikit/styles"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

func NewSecretView(
	ctx *context.Context,
	secret vault.Secret,
	asPlainText bool,
) components.TeaModel {
	container := ctx.InteractiveContainer
	v := vault.NewVault(ctx.Logger)
	var secretKeyCallbacks = []components.KeyCallback{
		{
			Key: "r", Label: "rename",
			Callback: func() error {
				in := components.TextInput{Key: "value", Prompt: "Enter the new secret name"}
				inputs, err := components.ProcessInputs(io.Theme(), &in)
				if err != nil {
					ctx.Logger.FatalErr(err)
				}
				newName := inputs.FindByKey("value").Value()
				if err := v.RenameSecret(secret.Reference, newName); err != nil {
					container.HandleError(fmt.Errorf("unable to rename secret: %w", err))
					return nil
				}
				LoadSecretListView(ctx, asPlainText)
				container.SetNotice("secret renamed", styles.NoticeLevelInfo)
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				in := components.TextInput{Key: "value", Prompt: "Enter the new secret value"}
				inputs, err := components.ProcessInputs(io.Theme(), &in)
				if err != nil {
					ctx.Logger.FatalErr(err)
				}
				newValue := inputs.FindByKey("value").Value()
				secretValue := vault.SecretValue(newValue)
				if err := v.SetSecret(secret.Reference, secretValue); err != nil {
					container.HandleError(fmt.Errorf("unable to edit secret: %w", err))
					return nil
				}
				LoadSecretListView(ctx, asPlainText)
				container.SetNotice("secret value updated", styles.NoticeLevelInfo)
				return nil
			},
		},
		{
			Key: "x", Label: "delete",
			Callback: func() error {
				if err := v.DeleteSecret(secret.Reference); err != nil {
					container.HandleError(fmt.Errorf("unable to delete secret: %w", err))
					return nil
				}
				LoadSecretListView(ctx, asPlainText)
				container.SetNotice("secret deleted", styles.NoticeLevelInfo)
				return nil
			},
		},
	}

	state := &components.TerminalState{
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewEntityView(state, &secret, components.FormatDocument, secretKeyCallbacks...)
}

func NewSecretListView(
	ctx *context.Context,
	secrets vault.SecretList,
	asPlainText bool,
) components.TeaModel {
	container := ctx.InteractiveContainer
	if len(secrets.Items()) == 0 {
		container.HandleError(fmt.Errorf("no secrets found"))
	}

	selectFunc := func(filterVal string) error {
		var secret vault.Secret
		var found bool
		for _, s := range secrets {
			if s.Reference == filterVal {
				secret = s
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("secret not found")
		}

		container.SetView(NewSecretView(ctx, secret, asPlainText))
		return nil
	}

	state := &components.TerminalState{
		Theme:  io.Theme(),
		Height: container.Height(),
		Width:  container.Width(),
	}
	return components.NewCollectionView(state, secrets, components.FormatList, selectFunc)
}

func LoadSecretListView(
	ctx *context.Context,
	asPlainText bool,
) {
	container := ctx.InteractiveContainer
	v := vault.NewVault(ctx.Logger)
	secrets, err := v.GetAllSecrets()
	if err != nil {
		ctx.Logger.FatalErr(err)
	}
	var secretList vault.SecretList
	for name, secret := range secrets {
		if asPlainText {
			secretList = append(secretList, vault.Secret{Reference: name, Secret: secret.PlainTextString()})
		} else {
			secretList = append(secretList, vault.Secret{Reference: name, Secret: secret.ObfuscatedString()})
		}
	}
	view := NewSecretListView(
		ctx,
		secretList,
		asPlainText,
	)
	container.SetView(view)
}
