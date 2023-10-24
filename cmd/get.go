package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/executable"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/io"
	configio "github.com/jahvon/flow/internal/io/config"
	executableio "github.com/jahvon/flow/internal/io/executable"
	workspaceio "github.com/jahvon/flow/internal/io/workspace"
	"github.com/jahvon/flow/internal/services/cache"
	"github.com/jahvon/flow/internal/workspace"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	GroupID: DataGroup.ID,
	Short:   "Print current flow data and metadata.",
}

var getConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Print the current flow config.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Panic().Msg("failed to load config")
		}

		outputFormatFlag, err := Flags.ValueFor(cmd, flags.OutputFormatFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		outputFormat, _ := outputFormatFlag.(string)
		configio.PrintRootConfig(io.OutputFormat(outputFormat), rootCfg)
	},
}

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

		workspaceFlag, err := Flags.ValueFor(cmd, flags.SpecificWorkspaceFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		workspaceName, _ := workspaceFlag.(string)
		if workspaceName != "" {
			if _, found := rootCfg.Workspaces[workspaceName]; !found {
				io.PrintErrorAndExit(fmt.Errorf("workspace %s not found", workspaceName))
			}
		} else {
			workspaceName = rootCfg.CurrentWorkspace
		}

		wsPath := rootCfg.Workspaces[workspaceName]
		wsCfg, err := workspace.LoadConfig(workspaceName, wsPath)
		if err != nil {
			log.Panic().Msgf("failed loading workspace config: %v", err)
		} else if wsCfg == nil {
			io.PrintErrorAndExit(fmt.Errorf("config not found for workspace %s", workspaceName))
		}

		outputFormatFlag, err := Flags.ValueFor(cmd, flags.OutputFormatFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		outputFormat, _ := outputFormatFlag.(string)
		workspaceio.PrintWorkspaceConfig(io.OutputFormat(outputFormat), wsCfg)
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

		outputFormatFlag, err := Flags.ValueFor(cmd, flags.OutputFormatFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		outputFormat, _ := outputFormatFlag.(string)

		tagsFlag, err := Flags.ValueFor(cmd, flags.FilterTagFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		tagsFilter, _ := tagsFlag.([]string)

		log.Debug().Msg("Loading workspace configs from cache")
		cacheData, err := cache.Get()
		if err != nil {
			log.Error().Err(err).Msg("failed to load workspace configs from cache")
		}

		if cacheData == nil {
			log.Debug().Msg("Cache data is nil; updating cache")
			cacheData, err = cache.Update()
			if err != nil || cacheData == nil {
				io.PrintErrorAndExit(fmt.Errorf("cache failure unrecoverable - %w", err))
			}
		}

		filteredWorkspaces := make([]workspace.Config, 0)
		for _, ws := range cacheData.Workspaces {
			if !ws.HasAnyTags(tagsFilter) {
				continue
			}
			filteredWorkspaces = append(filteredWorkspaces, *ws)
		}

		if len(filteredWorkspaces) == 0 {
			io.PrintErrorAndExit(fmt.Errorf("no workspaces found"))
		}
		workspaceio.PrintWorkspaceList(io.OutputFormat(outputFormat), filteredWorkspaces)
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

		executables, err := executable.FlagsToExecutableList(cmd, *Flags, rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		outputFormatFlag, err := Flags.ValueFor(cmd, flags.OutputFormatFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		outputFormat, _ := outputFormatFlag.(string)
		executableio.PrintExecutableList(io.OutputFormat(outputFormat), executables)
	},
}

var getExecutableCmd = &cobra.Command{
	Use:     "executable <identifier>",
	Aliases: []string{"exec"},
	Short:   "Find and print a discovered executable.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootCfg := config.LoadConfig()
		if rootCfg == nil {
			log.Panic().Msg("failed to load config")
		}

		agentFlag, err := Flags.ValueFor(cmd, flags.AgentTypeFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		agent, _ := agentFlag.(string)

		_, exec, err := executable.ArgsToExecutable(args, consts.AgentType(agent), rootCfg)
		if err != nil {
			io.PrintErrorAndExit(err)
		}

		outputFormatFlag, err := Flags.ValueFor(cmd, flags.OutputFormatFlag.Name)
		if err != nil {
			io.PrintErrorAndExit(err)
		}
		outputFormat, _ := outputFormatFlag.(string)
		executableio.PrintExecutable(io.OutputFormat(outputFormat), exec)
	},
}

func init() {
	registerFlagOrPanic(getConfigCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(getConfigCmd)

	registerFlagOrPanic(getWorkspaceCmd, *flags.OutputFormatFlag)
	registerFlagOrPanic(getWorkspaceCmd, *flags.SpecificWorkspaceFlag)
	getCmd.AddCommand(getWorkspaceCmd)

	registerFlagOrPanic(getWorkspacesCmd, *flags.OutputFormatFlag)
	registerFlagOrPanic(getWorkspacesCmd, *flags.FilterTagFlag)
	getCmd.AddCommand(getWorkspacesCmd)

	registerFlagOrPanic(getExecutablesCmd, *flags.OutputFormatFlag)
	registerFlagOrPanic(getExecutablesCmd, *flags.FilterAgentTypeFlag)
	registerFlagOrPanic(getExecutablesCmd, *flags.FilterTagFlag)
	registerFlagOrPanic(getExecutablesCmd, *flags.FilterNamespaceFlag)
	registerFlagOrPanic(getExecutablesCmd, *flags.ListWorkspaceContextFlag)
	registerFlagOrPanic(getExecutablesCmd, *flags.ListGlobalContextFlag)
	getExecutablesCmd.MarkFlagsMutuallyExclusive(
		flags.ListWorkspaceContextFlag.Name,
		flags.ListGlobalContextFlag.Name,
	)
	getCmd.AddCommand(getExecutablesCmd)

	registerFlagOrPanic(getExecutableCmd, *flags.OutputFormatFlag)
	registerFlagOrPanic(getExecutableCmd, *flags.AgentTypeFlag)
	getCmd.AddCommand(getExecutableCmd)

	rootCmd.AddCommand(getCmd)
}
