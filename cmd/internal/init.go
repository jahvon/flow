package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/crypto"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

func RegisterInitCmd(ctx *context.Context, rootCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize or restore the flow application state.",
	}
	registerInitConfigCmd(ctx, subCmd)
	registerInitWorkspaceCmd(ctx, subCmd)
	registerInitExecsCmd(ctx, subCmd)
	registerInitVaultCmd(ctx, subCmd)
	rootCmd.AddCommand(subCmd)
}

func registerInitConfigCmd(ctx *context.Context, initCmd *cobra.Command) {
	cfgCmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Initialize the flow global configuration.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { initConfigFunc(ctx, cmd, args) },
	}
	initCmd.AddCommand(cfgCmd)
}

func initConfigFunc(ctx *context.Context, _ *cobra.Command, _ []string) {
	logger := ctx.Logger
	inputs, err := components.ProcessInputs(io.Theme(), &components.TextInput{
		Key:    "confirm",
		Prompt: "This will overwrite your current flow configurations. Are you sure you want to continue? (y/n)",
	})
	if err != nil {
		logger.FatalErr(err)
	}
	resp := inputs.FindByKey("confirm").Value()
	if strings.ToLower(resp) != "y" && strings.ToLower(resp) != "yes" {
		logger.Warnf("Aborting", resp)
		return
	}

	if err := filesystem.InitUserConfig(); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess("Initialized flow global configurations")
}

func registerInitWorkspaceCmd(ctx *context.Context, initCmd *cobra.Command) {
	wsCmd := &cobra.Command{
		Use:     "workspace NAME PATH",
		Aliases: []string{"ws"},
		Short:   "Initialize and add a workspace to the list of known workspaces.",
		Args:    cobra.ExactArgs(2),
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { initWorkspaceFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, wsCmd, *flags.SetAfterCreateFlag)
	initCmd.AddCommand(wsCmd)
}

func initWorkspaceFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	name := args[0]
	path := args[1]

	userConfig := ctx.Config
	if _, found := userConfig.Workspaces[name]; found {
		logger.Fatalf("workspace %s already exists at %s", name, userConfig.Workspaces[name])
	}

	switch {
	case path == "":
		path = filepath.Join(filesystem.CachedDataDirPath(), name)
	case path == "." || strings.HasPrefix(path, "./"):
		wd, err := os.Getwd()
		if err != nil {
			logger.FatalErr(err)
		}
		if path == "." {
			path = wd
		} else {
			path = fmt.Sprintf("%s/%s", wd, path[2:])
		}
	case path == "~" || strings.HasPrefix(path, "~/"):
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

	if !filesystem.WorkspaceConfigExists(path) {
		if err := filesystem.InitWorkspaceConfig(name, path); err != nil {
			logger.FatalErr(err)
		}
	}
	userConfig.Workspaces[name] = path

	set := flags.ValueFor[bool](ctx, cmd, *flags.SetAfterCreateFlag, false)
	if set {
		userConfig.CurrentWorkspace = name
		logger.Infof("Workspace '%s' set as current workspace", name)
	}

	if err := filesystem.WriteUserConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}

	if err := cache.UpdateAll(logger); err != nil {
		logger.FatalErr(errors.Wrap(err, "failure updating cache"))
	}

	logger.PlainTextSuccess(fmt.Sprintf("Workspace '%s' created in %s", name, path))
}

func registerInitExecsCmd(ctx *context.Context, initCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "executables WORKSPACE_NAME DEFINITION_NAME [-p SUB_PATH] [-f FILE] [-t TEMPLATE]",
		Aliases: []string{"execs", "definitions", "defs"},
		Short:   "Add rendered executables from an executable definition template to a workspace",
		Long:    initExecLong,
		Args:    cobra.ExactArgs(2),
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { initExecFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, subCmd, *flags.SubPathFlag)
	RegisterFlag(ctx, subCmd, *flags.TemplateFlag)
	RegisterFlag(ctx, subCmd, *flags.FileFlag)
	MarkFlagMutuallyExclusive(subCmd, flags.TemplateFlag.Name, flags.FileFlag.Name)
	MarkOneFlagRequired(subCmd, flags.TemplateFlag.Name, flags.FileFlag.Name)
	MarkFlagFilename(ctx, subCmd, flags.FileFlag.Name)
	MarkFlagFilename(ctx, subCmd, flags.SubPathFlag.Name)
	initCmd.AddCommand(subCmd)
}

//nolint:gocognit
func initExecFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	workspaceName := args[0]
	definitionName := args[1]
	subPath := flags.ValueFor[string](ctx, cmd, *flags.SubPathFlag, false)
	template := flags.ValueFor[string](ctx, cmd, *flags.TemplateFlag, false)
	fileVal := flags.ValueFor[string](ctx, cmd, *flags.FileFlag, false)

	logger.Infof("Adding '%s' executables to '%s' workspace", definitionName, workspaceName)
	var flowFilePath string
	switch {
	case template == "" && fileVal == "":
		logger.Fatalf("one of -f or -t must be provided")
	case template != "" && fileVal != "":
		logger.Fatalf("only one of -f or -t can be provided")
	case template != "":
		if ctx.Config.Templates == nil {
			logger.Fatalf("template %s not found", template)
		}
		if path, found := ctx.Config.Templates[template]; !found {
			logger.Fatalf("template %s not found", template)
		} else {
			flowFilePath = path
		}
	case fileVal != "":
		if _, err := os.Stat(fileVal); os.IsNotExist(err) {
			logger.Fatalf("fileVal %s not found", fileVal)
		}
		flowFilePath = fileVal
	}

	execTemplate, err := filesystem.LoadFlowFileTemplate(flowFilePath)
	if err != nil {
		logger.FatalErr(err)
	}
	if err := execTemplate.Validate(); err != nil {
		logger.FatalErr(err)
	}
	execTemplate.SetContext(flowFilePath)

	wsPath, wsFound := ctx.Config.Workspaces[workspaceName]
	if !wsFound {
		logger.Fatalf("workspace %s not found", workspaceName)
	}
	ws, err := filesystem.LoadWorkspaceConfig(workspaceName, wsPath)
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
		inputs, err = components.ProcessInputs(io.Theme(), inputs...)
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

	if err := filesystem.InitExecutables(execTemplate, ws, definitionName, subPath); err != nil {
		logger.FatalErr(err)
	}

	logger.PlainTextSuccess(
		fmt.Sprintf(
			"Executables from %s added to %s\nPath: %s",
			definitionName, workspaceName, flowFilePath,
		))
}

func registerInitVaultCmd(ctx *context.Context, initCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:    "vault",
		Short:  "Create a new flow secret vault.",
		Args:   cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:    func(cmd *cobra.Command, args []string) { initVaultFunc(ctx, cmd, args) },
	}
	initCmd.AddCommand(subCmd)
}

func initVaultFunc(ctx *context.Context, _ *cobra.Command, _ []string) {
	logger := ctx.Logger
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
}

//nolint:lll
var initExecLong = `Add rendered executables from an executable definition template to a workspace.

The WORKSPACE_NAME is the name of the workspace to initialize the executables in.
The DEFINITION_NAME is the name of the definition to use in rendering the template.
This name will become the name of the file containing the copied executable definition.

One one of -f or -t must be provided and must point to a valid executable definition template.
The -p flag can be used to specify a sub-path within the workspace to create the executable definition and its artifacts.`
