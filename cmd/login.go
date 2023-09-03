package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/cmd/flags"
	"github.com/jahvon/tbox/internal/cmd/login"
	"github.com/jahvon/tbox/internal/config"
	"github.com/jahvon/tbox/internal/io"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to auth backend",
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}
		err := login.LoginWithFlags(cmd, rootCfg.Backends.Auth)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Successfully logged in")
	},
}

func init() {
	loginCmd.Flags().StringP(flags.PasswordFlagName, "p", "", "Password to use for login")
	loginCmd.Flags().StringP(flags.MasterKeyFlagName, "k", "", "Master key to use for login")
	loginCmd.MarkFlagsMutuallyExclusive(flags.PasswordFlagName, flags.MasterKeyFlagName)

	rootCmd.AddCommand(loginCmd)
}
