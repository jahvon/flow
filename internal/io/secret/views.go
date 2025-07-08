package secret

import (
	"fmt"

	"github.com/flowexec/tuikit"
	"github.com/flowexec/tuikit/themes"
	"github.com/flowexec/tuikit/types"
	"github.com/flowexec/tuikit/views"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/vault"
)

func NewSecretView(
	ctx *context.Context,
	secret vault.Secret,
	asPlainText bool,
) tuikit.View {
	container := ctx.TUIContainer
	v := vault.NewVault(ctx.Logger)
	var secretKeyCallbacks = []types.KeyCallback{
		{
			Key: "r", Label: "rename",
			Callback: func() error {
				form, err := views.NewFormView(
					container.RenderState(),
					&views.FormField{
						Key:   "value",
						Type:  views.PromptTypeText,
						Title: "Enter the new secret name",
					})
				if err != nil {
					container.HandleError(fmt.Errorf("encountered error creating the form: %w", err))
					return nil
				}
				if err := ctx.SetView(form); err != nil {
					container.HandleError(fmt.Errorf("unable to set view: %w", err))
					return nil
				}
				newName := form.FindByKey("value").Value()
				if err := v.RenameSecret(secret.Reference, newName); err != nil {
					container.HandleError(fmt.Errorf("unable to rename secret: %w", err))
					return nil
				}
				LoadSecretListView(ctx, asPlainText)
				container.SetNotice("secret renamed", themes.OutputLevelInfo)
				return nil
			},
		},
		{
			Key: "e", Label: "edit",
			Callback: func() error {
				form, err := views.NewFormView(
					container.RenderState(),
					&views.FormField{
						Key:   "value",
						Type:  views.PromptTypeMasked,
						Title: "Enter the new secret value",
					})
				if err != nil {
					container.HandleError(fmt.Errorf("encountered error creating the form: %w", err))
					return nil
				}
				if err := ctx.SetView(form); err != nil {
					container.HandleError(fmt.Errorf("unable to set view: %w", err))
					return nil
				}
				newValue := form.FindByKey("value").Value()
				secretValue := vault.SecretValue(newValue)
				if err := v.SetSecret(secret.Reference, secretValue); err != nil {
					container.HandleError(fmt.Errorf("unable to edit secret: %w", err))
					return nil
				}
				LoadSecretListView(ctx, asPlainText)
				container.SetNotice("secret value updated", themes.OutputLevelInfo)
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
				container.SetNotice("secret deleted", themes.OutputLevelInfo)
				return nil
			},
		},
	}

	return views.NewEntityView(container.RenderState(), &secret, types.EntityFormatDocument, secretKeyCallbacks...)
}

func NewSecretListView(
	ctx *context.Context,
	secrets vault.SecretList,
	asPlainText bool,
) tuikit.View {
	container := ctx.TUIContainer
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

		return container.SetView(NewSecretView(ctx, secret, asPlainText))
	}

	return views.NewCollectionView(container.RenderState(), secrets, types.CollectionFormatList, selectFunc)
}

func LoadSecretListView(
	ctx *context.Context,
	asPlainText bool,
) {
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
	if err := ctx.SetView(view); err != nil {
		ctx.Logger.FatalErr(err)
	}
}
