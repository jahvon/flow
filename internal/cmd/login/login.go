package login

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/backend"
	"github.com/jahvon/tbox/internal/backend/consts"
	"github.com/jahvon/tbox/internal/backend/conv"
	"github.com/jahvon/tbox/internal/cmd/flags"
	"github.com/jahvon/tbox/internal/io"
)

const AuthEnvVar = "TBOX_MASTER_KEY"

var log = io.Log()

func LoginWithFlags(cmd *cobra.Command, authCfg *backend.AuthConfig) error {
	password := cmd.Flag(flags.PasswordFlagName).Value.String()
	masterKey := cmd.Flag(flags.MasterKeyFlagName).Value.String()

	if authCfg == nil {
		return errors.New("auth config is not set")
	}

	preferredAuth := authCfg.PreferredMode
	modeNotSpecified := password == "" && masterKey == ""
	if modeNotSpecified && preferredAuth == consts.ModeMasterKey {
		envMasterKey, ok := os.LookupEnv(AuthEnvVar)
		if ok && envMasterKey != "" {
			log.Debug().Msg("Using master key from environment variable")
			masterKey = envMasterKey
		}

		promptMasterKey := io.AskForMasterKey()
		if promptMasterKey != "" {
			log.Debug().Msg("Using master key from prompt")
			masterKey = promptMasterKey
		}

		if masterKey == "" {
			return errors.New("master key not found")
		}
	} else if modeNotSpecified && preferredAuth == consts.ModePassword {
		promptPassword := io.AskForPassword()
		if promptPassword != "" {
			log.Debug().Msg("Using password from prompt")
			password = promptPassword
		}

		if password == "" {
			return errors.New("password not found")
		}
	} else if modeNotSpecified {
		return errors.New("preferred auth mode not set; please set either password or master key")
	}

	if password != "" && masterKey != "" {
		log.Info().Msg("password and master key provided, using master key")
	}

	authBackend, err := conv.AuthBackendFromName(authCfg.Backend)
	if err != nil {
		return err
	}

	if masterKey != "" {
		if err := authCfg.LoginWithMasterKey(authBackend, masterKey); err != nil {
			return err
		}
		log.Info().Msg("successfully logged in with master key")
		return nil
	} else if password != "" {
		if err := authCfg.LoginWithPassword(authBackend, password); err != nil {
			return err
		}
		log.Info().Msg("successfully logged in with password")
		return nil
	} else {
		return errors.New("password or master key not provided")
	}
}
