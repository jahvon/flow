package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
	workspaceio "github.com/jahvon/flow/internal/io/workspace"
)

var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w", "ws"},
	Short:   "Manage flow workspaces.",
}

var workspaceSetCmd = &cobra.Command{
	Use:     "set <name>",
	Aliases: []string{"s"},
	Short:   "Change the current workspace.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace := args[0]
		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		if _, found := userConfig.Workspaces[workspace]; !found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s not found", workspace))
		}
		userConfig.CurrentWorkspace = workspace

		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess("Workspace set to " + workspace)
	},
}

var workspaceGetCmd = &cobra.Command{
	Use:     "get <name>",
	Aliases: []string{"g"},
	Short:   "Get a workspace's configuration. If the name is omitted, the current workspace is used.",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		var workspaceName string
		if len(args) == 1 {
			workspaceName = args[0]
		} else {
			workspaceName = userConfig.CurrentWorkspace
		}

		if _, found := userConfig.Workspaces[workspaceName]; !found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s not found", workspaceName))
		}

		wsPath := userConfig.Workspaces[workspaceName]
		wsCfg, err := file.LoadWorkspaceConfig(workspaceName, wsPath)
		if err != nil {
			log.Panic().Msgf("failed loading workspace config: %v", err)
		} else if wsCfg == nil {
			io.PrintErrorAndExit(fmt.Errorf("config not found for workspace %s", workspaceName))
		}

		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		workspaceio.PrintWorkspaceConfig(io.OutputFormat(outputFormat), wsCfg)
	},
}

var workspaceAddCmd = &cobra.Command{
	Use:     "add <name> <path>",
	Aliases: []string{"a"},
	Short:   "Add a new workspace.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := args[1]

		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		if _, found := userConfig.Workspaces[name]; found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s already exists at %s", name, userConfig.Workspaces[name]))
		}

		if path == "" {
			path = file.CachedDataDirPath()
		} else if path == "." || strings.HasPrefix(path, "./") {
			wd, err := os.Getwd()
			if err != nil {
				io.PrintErrorAndExit(err)
			}
			if path == "." {
				path = wd
			} else {
				path = fmt.Sprintf("%s/%s", wd, path[2:])
			}
		} else if path == "~" || strings.HasPrefix(path, "~/") {
			hd, err := os.UserHomeDir()
			if err != nil {
				io.PrintErrorAndExit(err)
			}
			if path == "~" {
				path = hd
			} else {
				path = fmt.Sprintf("%s/%s", hd, path[2:])
			}
		}

		if err := file.InitWorkspaceConfig(name, path); err != nil {
			io.PrintErrorAndExit(err)
		}
		io.PrintSuccess(fmt.Sprintf("Workspace %s created in %s", name, path))

		set := getFlagValue[bool](cmd, *flags.SetAfterCreateFlag)
		if set {
			userConfig.CurrentWorkspace = name
			if err := file.WriteUserConfig(userConfig); err != nil {
				io.PrintErrorAndExit(err)
			}
			io.PrintInfo(fmt.Sprintf("Workspace %s set as current workspace", name))
		}

		wsCache := cache.NewWorkspaceCache()
		if err := wsCache.Update(); err != nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to update cache - %w", err))
		} else {
			io.PrintInfo("Workspace cache updated")
		}
	},
}

var workspaceList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "ParameterList workspace configurations.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := newCtx(cmd)
		outputFormat := getFlagValue[string](cmd, *flags.OutputFormatFlag)
		tagsFilter := getFlagValue[config.Tags](cmd, *flags.FilterTagFlag)

		log.Debug().Msg("Loading workspace configs from cache")
		workspaceCache, err := ctx.WorkspacesCache.Get()
		if err != nil {
			log.Error().Err(err).Msg("failed to load workspace configs from cache")
		}

		filteredWorkspaces := make([]config.WorkspaceConfig, 0)
		for _, ws := range workspaceCache.Workspaces {
			if !ws.Tags.HasAnyTag(tagsFilter) {
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

var workspaceRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm", "r"},
	Short:   "Delete an existing workspace.",
	Long: "Delete an existing workspace. File contents will remain in the corresponding directory but the " +
		"workspace will be unlinked from the flow user configurations.\nNote: You cannot delete the current workspace.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		confirmed := io.AskYesNo("Are you sure you want to delete the workspace '" + name + "'?")
		if !confirmed {
			io.PrintWarning("Aborting")
			return
		}

		userConfig := file.LoadUserConfig()
		if userConfig == nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to load user config"))
		}
		if err := userConfig.Validate(); err != nil {
			io.PrintErrorAndExit(err)
		}

		if name == userConfig.CurrentWorkspace {
			io.PrintErrorAndExit(fmt.Errorf("cannot delete the current workspace"))
		}
		if _, found := userConfig.Workspaces[name]; !found {
			io.PrintErrorAndExit(fmt.Errorf("workspace %s was not found", name))
		}

		delete(userConfig.Workspaces, name)
		if err := file.WriteUserConfig(userConfig); err != nil {
			io.PrintErrorAndExit(err)
		}

		io.PrintWarning(fmt.Sprintf("Workspace %s deleted", name))

		wsCache := cache.NewWorkspaceCache()
		if err := wsCache.Update(); err != nil {
			io.PrintErrorAndExit(fmt.Errorf("failed to update cache - %w", err))
		} else {
			io.PrintInfo("Workspace cache updated")
		}
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceSetCmd)

	registerFlagOrPanic(workspaceGetCmd, *flags.OutputFormatFlag)
	workspaceCmd.AddCommand(workspaceGetCmd)

	registerFlagOrPanic(workspaceAddCmd, *flags.SetAfterCreateFlag)
	workspaceCmd.AddCommand(workspaceAddCmd)

	registerFlagOrPanic(workspaceList, *flags.OutputFormatFlag)
	registerFlagOrPanic(workspaceList, *flags.FilterTagFlag)
	workspaceCmd.AddCommand(workspaceList)

	workspaceCmd.AddCommand(workspaceRemoveCmd)

	rootCmd.AddCommand(workspaceCmd)
}
