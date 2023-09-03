package config

import (
	"fmt"

	"github.com/jahvon/tbox/internal/backend"
	"github.com/jahvon/tbox/internal/backend/consts"
	"github.com/jahvon/tbox/internal/backend/conv"
	"github.com/jahvon/tbox/internal/io"
)

func SetSecretBackend(config *RootConfig, secretConfig *backend.SecretConfig) error {
	if secretConfig == nil {
		log.Info().Msg("secret config not changed")
		return nil
	}
	backendName := secretConfig.Backend
	curBackend := config.Backends.Secret.Backend
	io.PrintWarning("Updating the auth backend may break functionality of the current backend!")
	if backendName != curBackend {
		yes := io.AskYesNo(
			fmt.Sprintf(
				"Are you sure you want to change the secret backend from %s to %s?",
				curBackend,
				backendName,
			),
		)
		if !yes {
			return nil
		}
	} else {
		log.Info().Msg("secret backend not changed")
		return nil
	}

	secretBackend, err := conv.SecretBackendFromName(backendName)
	if err != nil {
		return err
	}
	if err := secretBackend.InitializeBackend(); err != nil {
		return fmt.Errorf("unable to initialize secret backend - %v", err)
	}

	config.Backends.Secret.Backend = backendName
	return writeConfigFile(config)
}

func SetAuthBackend(config *RootConfig, authConfig *backend.AuthConfig) error {
	if authConfig == nil {
		log.Info().Msg("auth config not changed")
		return nil
	}
	backendName := authConfig.Backend
	curBackend := config.Backends.Auth.Backend
	io.PrintWarning("Updating the auth backend may break functionality of the current backend!")
	if backendName != curBackend {
		yes := io.AskYesNo(
			fmt.Sprintf(
				"Are you sure you want to change the auth backend from '%s' to '%s'?",
				curBackend,
				backendName,
			),
		)
		if !yes {
			return nil
		}
	}

	authBackend, err := conv.AuthBackendFromName(backendName)
	if err != nil {
		return err
	}

	modeChanged := config.Backends.Auth.PreferredMode != authConfig.PreferredMode

	var mk string
	var newMK bool
	if authConfig.PreferredMode == "" || authConfig.PreferredMode == consts.ModeMasterKey {
		if !modeChanged {
			yes := io.AskYesNo("Would you like to generate a new master key?")
			if yes {
				newMK = true
			}
		} else {
			newMK = true
		}

		if newMK {
			log.Trace().Msg("generating a new master key")
			var err error
			mk, err = authConfig.InitializeWithRandomMasterKey(authBackend)
			if err != nil {
				return fmt.Errorf("unable to initialize auth config with a new master key - %v", err)
			}
		}
	} else if authConfig.PreferredMode == consts.ModePassword {
		if !modeChanged {
			yes := io.AskYesNo("Would you like to change your password?")
			if yes {
				newMK = true
			}
		} else {
			newMK = true
		}

		if newMK {
			var err error
			passwordInput := io.AskForPassword()
			if passwordInput == "" {
				return fmt.Errorf("password cannot be empty")
			}

			mk, err = authConfig.InitializeWithPassword(authBackend, passwordInput)
			if err != nil {
				return fmt.Errorf("unable to initialize auth config with password - %v", err)
			}
		}
	} else {
		return fmt.Errorf("unknown preferred mode - %v", authConfig.PreferredMode)
	}

	if newMK {
		io.PrintNotice(fmt.Sprintf("Your master key is: %s\nSave this somewhere safe!", mk))
	}

	config.Backends.Auth = authConfig
	return writeConfigFile(config)
}
