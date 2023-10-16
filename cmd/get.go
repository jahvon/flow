package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/executable"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/io"
	executableio "github.com/jahvon/flow/internal/io/executable"
)

// getCmd represents the get command.
var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	GroupID: CrudGroup.ID,
	Short:   "Get the current value of a configuration, environment, or workspace option.",
}

// getWorkspaceCmd represents the get workspace subcommand.
var getWorkspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w"},
	Short:   "Print the current workspace.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Panic().Msg("failed to load config")
		}
		wsPath := rootCfg.Workspaces[rootCfg.CurrentWorkspace]
		wsInfo := fmt.Sprintf("%s (%s)", rootCfg.CurrentWorkspace, wsPath)
		io.PrintNotice(wsInfo)
	},
}

var getWorkspacesCmd = &cobra.Command{
	Use:     "workspaces",
	Aliases: []string{"wss"},
	Short:   "Print a list of discovered workspaces.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Panic().Msg("failed to load config")
		}

		log.Info().Msgf("Printing %d workspaces", len(rootCfg.Workspaces))
		tableRows := pterm.TableData{{"Name", "Location"}}
		for ws, wsPath := range rootCfg.Workspaces {
			tableRows = append(tableRows, []string{ws, wsPath})
		}
		err := pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableRows).Render()
		if err != nil {
			log.Panic().Msgf("Failed to render workspace list - %v", err)
		}
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
			log.Panic().Msg("failed to load config")
		}

		executables, err := executable.FlagsToExecutableList(cmd, rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		outputFormat := cmd.Flag(flags.OutputFormatFlagName).Value.String()
		executableio.PrintExecutableList(io.OutputFormat(outputFormat), executables)
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
			log.Panic().Msg("failed to load config")
		}

		agent := cmd.Flag(flags.AgentTypeFlagName)
		if agent == nil || !agent.Changed {
			log.Panic().Msg("agent type is required")
		}

		_, exec, err := executable.ArgsToExecutable(args, consts.AgentType(agent.Value.String()), rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		outputFormat := cmd.Flag(flags.OutputFormatFlagName).Value.String()
		executableio.PrintExecutable(io.OutputFormat(outputFormat), exec)
	},
}

func init() {
	getCmd.AddCommand(getWorkspaceCmd)
	getCmd.AddCommand(getWorkspacesCmd)

	getExecutablesCmd.Flags().StringP(
		flags.OutputFormatFlagName,
		"o",
		"default",
		"Output format. One of: default, yaml, json, jsonp.",
	)
	getExecutablesCmd.Flags().StringP(
		flags.AgentTypeFlagName,
		"a",
		"",
		"Filter executables by agent type.",
	)
	getExecutablesCmd.Flags().StringP(
		flags.TagFlagName,
		"t",
		"",
		"Filter executables by tag.",
	)
	getExecutablesCmd.Flags().StringP(
		flags.NamespaceFlagName,
		"n",
		"",
		"Filter executables by namespace.",
	)
	getExecutablesCmd.Flags().StringP(
		flags.WorkspaceContextFlagName, "w", "", "Filter executables by workspace.",
	)
	getExecutablesCmd.Flags().BoolP(
		flags.GlobalContextFlagName, "g", false, "List executables across all workspaces.",
	)
	getExecutablesCmd.MarkFlagsMutuallyExclusive(flags.WorkspaceContextFlagName, flags.GlobalContextFlagName)
	getCmd.AddCommand(getExecutablesCmd)

	getExecutableCmd.Flags().StringP(
		flags.OutputFormatFlagName,
		"o",
		"default",
		"Output format. One of: default, yaml, json, jsonp.",
	)
	getExecutableCmd.Flags().StringP(
		flags.AgentTypeFlagName,
		"a",
		"",
		fmt.Sprintf("Executable agent type. One of: %s", consts.ValidAgentTypes),
	)
	_ = getExecutableCmd.MarkFlagRequired(flags.AgentTypeFlagName)
	getCmd.AddCommand(getExecutableCmd)

	rootCmd.AddCommand(getCmd)
}
