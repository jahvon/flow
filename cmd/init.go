package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/crypto"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or restore the flow application state.",
}

var configInitCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Initialize the flow global configuration.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		inputs, err := components.ProcessInputs(io.Styles(), &components.TextInput{
			Key:    "confirm",
			Prompt: "This will overwrite your current flow configurations. Are you sure you want to continue? (y/n)",
		})
		if err != nil {
			logger.FatalErr(err)
		}
		resp := inputs.FindByKey("confirm").Value()
		if truthy, _ := strconv.ParseBool(resp); !truthy {
			logger.Warnf("Aborting")
			return
		}

		if err := file.InitUserConfig(); err != nil {
			logger.FatalErr(err)
		}
		logger.PlainTextSuccess("Initialized flow global configurations")
	},
}

var workspaceInitCmd = &cobra.Command{
	Use:     "workspace <name> <path>",
	Aliases: []string{"ws"},
	Short:   "Initialize and add a workspace to the list of known workspaces.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		name := args[0]
		path := args[1]

		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		userConfig := curCtx.UserConfig
		if _, found := userConfig.Workspaces[name]; found {
			logger.Fatalf("workspace %s already exists at %s", name, userConfig.Workspaces[name])
		}

		if path == "" {
			path = filepath.Join(file.CachedDataDirPath(), name)
		} else if path == "." || strings.HasPrefix(path, "./") {
			wd, err := os.Getwd()
			if err != nil {
				logger.FatalErr(err)
			}
			if path == "." {
				path = wd
			} else {
				path = fmt.Sprintf("%s/%s", wd, path[2:])
			}
		} else if path == "~" || strings.HasPrefix(path, "~/") {
			hd, err := os.UserHomeDir()
			if err != nil {
				logger.FatalErr(err)
			}
			if path == "~" {
				path = hd
			} else {
				path = fmt.Sprintf("%s/%s", hd, path[2:])
			}
		}

		if err := file.InitWorkspaceConfig(name, path); err != nil {
			logger.FatalErr(err)
		}
		userConfig.Workspaces[name] = path

		set := getFlagValue[bool](cmd, *flags.SetAfterCreateFlag)
		if set {
			userConfig.CurrentWorkspace = name
			logger.Infof("Workspace '%s' set as current workspace", name)
		}

		if err := file.WriteUserConfig(userConfig); err != nil {
			logger.FatalErr(err)
		}

		if err := cache.UpdateAll(logger); err != nil {
			logger.FatalErr(errors.Wrap(err, "failure updating cache"))
		}

		logger.PlainTextSuccess(fmt.Sprintf("Workspace '%s' created in %s", name, path))
	},
}

var execsInitCmd = &cobra.Command{
	Use:     "executables WORKSPACE_NAME DEFINITION_NAME [-p SUB_PATH] [-f FILE] [-t TEMPLATE]",
	Aliases: []string{"execs", "definitions", "defs"},
	Short:   "Add rendered executables from an executable definition template to a workspace",
	//nolint:lll
	Long: `Add rendered executables from an executable definition template to a workspace.

The WORKSPACE_NAME is the name of the workspace to initialize the executables in.
The DEFINITION_NAME is the name of the definition to use in rendering the template. 
This name will become the name of the file containing the copied executable definition.

One one of -f or -t must be provided and must point to a valid executable definition template.
The -p flag can be used to specify a sub-path within the workspace to create the executable definition and its artifacts.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		workspaceName := args[0]
		definitionName := args[1]
		subPath := getFlagValue[string](cmd, *flags.SubPathFlag)

		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Notice = fmt.Sprintf("Adding '%s' executables to '%s' workspace", definitionName, workspaceName)
			header.Print()
		}

		templateFlag := getFlagValue[string](cmd, *flags.TemplateFlag)
		fileFlag := getFlagValue[string](cmd, *flags.FileFlag)
		var definitionPath string
		switch {
		case templateFlag == "" && fileFlag == "":
			logger.Fatalf("one of -f or -t must be provided")
		case templateFlag != "" && fileFlag != "":
			logger.Fatalf("only one of -f or -t can be provided")
		case templateFlag != "":
			if curCtx.UserConfig.Templates == nil {
				logger.Fatalf("template %s not found", templateFlag)
			}
			if path, found := curCtx.UserConfig.Templates[templateFlag]; !found {
				logger.Fatalf("template %s not found", templateFlag)
			} else if found {
				definitionPath = path
			}
		case fileFlag != "":
			if _, err := os.Stat(fileFlag); os.IsNotExist(err) {
				logger.Fatalf("file %s not found", fileFlag)
			}
			definitionPath = fileFlag
		}
		execTemplate, err := file.LoadExecutableDefinitionTemplate(definitionPath)
		if err != nil {
			logger.FatalErr(err)
		}
		if err := execTemplate.Validate(); err != nil {
			logger.FatalErr(err)
		}
		execTemplate.SetContext(definitionPath)

		wsPath, wsFound := curCtx.UserConfig.Workspaces[workspaceName]
		if !wsFound {
			logger.Fatalf("workspace %s not found", workspaceName)
		}
		ws, err := file.LoadWorkspaceConfig(workspaceName, wsPath)
		if err != nil {
			logger.FatalErr(err)
		}
		ws.SetContext(workspaceName, wsPath)

		if len(execTemplate.Data) != 0 {
			var inputs []*components.TextInput
			for _, entry := range execTemplate.Data {
				inputs = append(inputs, &components.TextInput{
					Key:         entry.Key,
					Prompt:      entry.Prompt,
					Placeholder: entry.Default,
				})
			}
			inputs, err = components.ProcessInputs(io.Styles(), inputs...)
			if err != nil {
				logger.FatalErr(err)
			}
			for _, input := range inputs {
				execTemplate.Data.Set(input.Key, input.Value())
			}
			if err := execTemplate.Data.ValidateValues(); err != nil {
				logger.FatalErr(err)
			}
		}

		if err := file.InitExecutables(execTemplate, ws, definitionName, subPath); err != nil {
			logger.FatalErr(err)
		}

		logger.PlainTextSuccess(
			fmt.Sprintf(
				"Executables from %s added to %s\nPath: %s",
				definitionName, workspaceName, definitionPath,
			))
	},
}

var vaultInitCmd = &cobra.Command{
	Use:   "vault",
	Short: "Create a new flow secret vault.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger := curCtx.Logger
		if interactiveUIEnabled() {
			header := headerForCurCtx()
			header.Print()
		}
		generatedKey, err := crypto.GenerateKey()
		if err != nil {
			logger.FatalErr(err)
		}
		if err = vault.RegisterEncryptionKey(generatedKey); err != nil {
			logger.FatalErr(err)
		}

		logger.PlainTextSuccess(fmt.Sprintf("Your vault encryption key is: %s", generatedKey))
		newKeyMsg := fmt.Sprintf(
			"You will need this key to modify your vault data. Store it somewhere safe!\n"+
				"Set this value to the %s environment variable if you do not want to be prompted for it every time.",
			vault.EncryptionKeyEnvVar,
		)
		logger.PlainTextInfo(newKeyMsg)
	},
}

func init() {
	initCmd.AddCommand(configInitCmd)

	registerFlagOrPanic(workspaceInitCmd, *flags.SetAfterCreateFlag)
	initCmd.AddCommand(workspaceInitCmd)

	registerFlagOrPanic(execsInitCmd, *flags.SubPathFlag)
	registerFlagOrPanic(execsInitCmd, *flags.TemplateFlag)
	registerFlagOrPanic(execsInitCmd, *flags.FileFlag)
	execsInitCmd.MarkFlagsMutuallyExclusive(flags.TemplateFlag.Name, flags.FileFlag.Name)
	execsInitCmd.MarkFlagsOneRequired(flags.TemplateFlag.Name, flags.FileFlag.Name)
	_ = execsInitCmd.MarkFlagFilename(flags.FileFlag.Name)
	_ = execsInitCmd.MarkFlagFilename(flags.SubPathFlag.Name)
	initCmd.AddCommand(execsInitCmd)

	initCmd.AddCommand(vaultInitCmd)

	rootCmd.AddCommand(initCmd)
}
