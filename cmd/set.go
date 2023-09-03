package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/tbox/internal/backend/consts"
	"github.com/jahvon/tbox/internal/cmd/flags"
	"github.com/jahvon/tbox/internal/cmd/set"
	"github.com/jahvon/tbox/internal/cmd/utils"
	"github.com/jahvon/tbox/internal/config"
	"github.com/jahvon/tbox/internal/io"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Update a configuration, environment, or workspace option.",
}

// setWorkspaceCmd represents the set workspace subcommand
var setWorkspaceCmd = &cobra.Command{
	Use:     "workspace <name>",
	Aliases: []string{"w"},
	Short:   "Change the current workspace.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace := args[0]
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		if err := config.SetCurrentWorkspace(rootCfg, workspace); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Workspace set to " + workspace)
	},
}

// setAuthCmd represents the set auth subcommand
var setAuthCmd = &cobra.Command{
	Use:     "auth",
	Aliases: []string{"a"},
	Short:   "Options for updating auth config and data.",
}

// setAuthBackendCmd represents the set backend auth subcommand
var setAuthBackendCmd = &cobra.Command{
	Use:     "backend",
	Aliases: []string{"b"},
	Short:   "Update the authentication backend configurations.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}
		authConfig, err := set.FlagsToAuthConfig(cmd, rootCfg.Backends.Auth)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		if err := config.SetAuthBackend(rootCfg, authConfig); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintSuccess("Authentication backend configurations updated")
	},
}

// setSecretCmd represents the set secret subcommand
var setSecretCmd = &cobra.Command{
	Use:     "secret",
	Aliases: []string{"s"},
	Short:   "Options for updating secret config and data.",
}

// setSecretBackendCmd represents the set secret subcommand
var setSecretBackendCmd = &cobra.Command{
	Use:     "backend",
	Aliases: []string{"b"},
	Short:   "Update the secret backend configurations.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}
		secretConfig, err := set.FlagsToSecretConfig(cmd, rootCfg.Backends.Secret)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		if err := config.SetSecretBackend(rootCfg, secretConfig); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintSuccess("Secret backend configurations updated")
	},
}

// setSecretDataCmd represents the set secret data subcommand
var setSecretDataCmd = &cobra.Command{
	Use:     "data <key> <value>",
	Aliases: []string{"d"},
	Short:   "Create or update a secret.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		key, value, err := set.ArgsToSecretKV(args)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		context, err := utils.ValidateAndGetContext(cmd, rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		if err := rootCfg.Backends.Secret.CurrentBackend().SetSecret(context, key, value); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintSuccess(fmt.Sprintf("Secret '%s' set within '%s' context", key, context))
	},
}

// setParamCmd represents the set param subcommand
var setParamCmd = &cobra.Command{
	Use:     "param <key>",
	Aliases: []string{"p"},
	Short:   "Create or update a parameter.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		// todo
		// context, err := utils.ValidateAndGetContext(cmd, rootCfg)
		// if err != nil {
		// 	io.PrintErrorAndExit(err)
		// }
		//
		// if context == "global" {
		// 	param, err := set.FlagsToParameter(cmd)
		// 	if err != nil {
		// 		io.PrintErrorAndExit(err)
		// 	}
		// 	io.PrintSuccess(fmt.Sprintf("Parameter '%s' set with global scope", key))
		// } else {
		// 	for ws, wsPath := range currentConfig.Workspaces {
		// 		if ws == context {
		//
		// 			io.PrintSuccess(fmt.Sprintf("Parameter '%s' set in workspace '%s' config", key, currentWorkspaceName))
		// 		}
		// 	}
		// }
		log.Panic().Msg("not implemented")
	},
}

func init() {
	// Auth sub commands
	setAuthBackendCmd.Flags().StringP(
		flags.BackendFlagName,
		"n",
		"",
		fmt.Sprintf("Set the authentication backend to use. Valid options are: %s", consts.AuthBackends),
	)
	setAuthBackendCmd.Flags().StringP(
		flags.PreferredModeFlagName,
		"m",
		"",
		fmt.Sprintf("Set the preferred authentication mode to use. Valid options are: %s", consts.AuthModes),
	)
	setAuthBackendCmd.Flags().BoolP(
		flags.RememberMeFlagName,
		"r",
		false,
		"Save auth session for a given duration.",
	)
	setAuthBackendCmd.Flags().DurationP(
		flags.RememberMeDurationFlagName,
		"d",
		0,
		"Set the duration to remember the user",
	)
	setAuthCmd.AddCommand(setAuthBackendCmd)
	setCmd.AddCommand(setAuthCmd)

	// Param sub commands
	setParamCmd.Flags().StringP(
		flags.WorkspaceContextFlagName,
		"w",
		"",
		"Workspace to set the parameter in, defaults to current workspace if not set",
	)
	setParamCmd.Flags().BoolP(
		flags.GlobalContextFlagName,
		"g",
		false,
		"Set the parameter to be globally accessible",
	)
	setParamCmd.Flags().StringP(
		flags.TextValueFlagName,
		"t",
		"",
		"Set the parameter value to the specified text",
	)
	setParamCmd.Flags().StringP(
		flags.SecretRefFlagName,
		"r",
		"",
		"Set the parameter value to the specified secret reference (key)",
	)
	setParamCmd.MarkFlagsMutuallyExclusive(flags.TextValueFlagName, flags.SecretRefFlagName)
	setCmd.AddCommand(setParamCmd)

	// Secret sub commands
	setSecretBackendCmd.Flags().StringP(
		flags.BackendFlagName,
		"n",
		"",
		fmt.Sprintf("Set the secret backend to use. Valid options are: %s", consts.SecretBackends),
	)
	setSecretDataCmd.Flags().StringP(
		flags.WorkspaceContextFlagName,
		"w",
		"",
		"Workspace to set the parameter in, defaults to current workspace if not set",
	)
	setSecretDataCmd.Flags().BoolP(
		flags.GlobalContextFlagName,
		"g",
		false,
		"Set the parameter to be globally accessible",
	)
	setSecretDataCmd.MarkFlagsMutuallyExclusive("workspace", "global")
	setSecretCmd.AddCommand(setSecretBackendCmd)
	setSecretCmd.AddCommand(setSecretDataCmd)
	setCmd.AddCommand(setSecretCmd)

	setCmd.AddCommand(setWorkspaceCmd)

	rootCmd.AddCommand(setCmd)
}
