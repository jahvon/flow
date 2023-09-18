package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/executable"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/io"
	executable2 "github.com/jahvon/flow/internal/io/executable"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	GroupID: CrudGroup.ID,
	Short:   "Get the current value of a configuration, environment, or workspace option.",
}

// getWorkspaceCmd represents the get workspace subcommand
var getWorkspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w"},
	Short:   "Print the current workspace.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}
		io.PrintNotice(rootCfg.CurrentWorkspace)
	},
}

var getExecutablesCmd = &cobra.Command{
	Use:     "executables",
	Aliases: []string{"execs"},
	Short:   "Print a list of discovered executables.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		executables, err := executable.FlagsToExecutableList(cmd, rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		outputFormat := cmd.Flag(flags.OutputFormatFlagName).Value.String()
		executable2.PrintExecutableList(io.OutputFormat(outputFormat), executables)
	},
}

var getExecutableCmd = &cobra.Command{
	Use:     "executable",
	Aliases: []string{"exec"},
	Short:   "Find and print a discovered executable.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Fatal().Msg("failed to load config")
		}

		agent := cmd.Flag(flags.AgentTypeFlagName)
		if agent == nil || !agent.Changed {
			log.Fatal().Msg("agent type is required")
		}

		_, exec, err := executable.ArgsToExecutable(args, consts.AgentType(agent.Value.String()), rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		outputFormat := cmd.Flag(flags.OutputFormatFlagName).Value.String()
		executable2.PrintExecutable(io.OutputFormat(outputFormat), exec)
	},
}

func init() {
	getCmd.AddCommand(getWorkspaceCmd)

	getExecutablesCmd.Flags().StringP(flags.OutputFormatFlagName, "o", "default", "Output format. One of: default, yaml, json, jsonp.")
	getExecutablesCmd.Flags().StringP(flags.AgentTypeFlagName, "a", "", "Filter executables by agent type.")
	getExecutablesCmd.Flags().StringP(flags.TagFlagName, "t", "", "Filter executables by tag.")
	getExecutablesCmd.Flags().StringP(flags.NamespaceFlagName, "n", "", "Filter executables by namespace.")
	getExecutablesCmd.Flags().StringP(
		flags.WorkspaceContextFlagName, "w", "", "Filter executables by workspace.",
	)
	getExecutablesCmd.Flags().BoolP(
		flags.GlobalContextFlagName, "g", false, "List executables across all workspaces.",
	)
	getExecutablesCmd.MarkFlagsMutuallyExclusive(flags.WorkspaceContextFlagName, flags.GlobalContextFlagName)
	getCmd.AddCommand(getExecutablesCmd)

	getExecutableCmd.Flags().StringP(flags.OutputFormatFlagName, "o", "default", "Output format. One of: default, yaml, json, jsonp.")
	getExecutableCmd.Flags().StringP(flags.AgentTypeFlagName, "a", "", fmt.Sprintf("Executable agent type. One of: %s", consts.ValidAgentTypes))
	_ = getExecutableCmd.MarkFlagRequired(flags.AgentTypeFlagName)
	getCmd.AddCommand(getExecutableCmd)

	rootCmd.AddCommand(getCmd)

}
