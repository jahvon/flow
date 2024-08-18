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
	"github.com/jahvon/flow/internal/templates"
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

	if err := filesystem.InitConfig(); err != nil {
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

	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}

	if err := cache.UpdateAll(logger); err != nil {
		logger.FatalErr(errors.Wrap(err, "failure updating cache"))
	}

	logger.PlainTextSuccess(fmt.Sprintf("Workspace '%s' created in %s", name, path))
}

func registerInitExecsCmd(ctx *context.Context, initCmd *cobra.Command) {
	subCmd := &cobra.Command{
		Use:     "executables FLOWFILE_NAME [-w WORKSPACE ] [-o OUTPUT_DIR] [-f FILE | -t TEMPLATE]",
		Aliases: []string{"execs", "flowfile"},
		Short:   "Add rendered executables from an executable definition template to a workspace",
		Long:    initExecLong,
		Args:    cobra.MaximumNArgs(1),
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { initExecFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, subCmd, *flags.TemplateOutputPathFlag)
	RegisterFlag(ctx, subCmd, *flags.TemplateFlag)
	RegisterFlag(ctx, subCmd, *flags.TemplateFilePathFlag)
	RegisterFlag(ctx, subCmd, *flags.TemplateWorkspaceFlag)
	MarkFlagMutuallyExclusive(subCmd, flags.TemplateFlag.Name, flags.TemplateFilePathFlag.Name)
	MarkOneFlagRequired(subCmd, flags.TemplateFlag.Name, flags.TemplateFilePathFlag.Name)
	MarkFlagFilename(ctx, subCmd, flags.TemplateFilePathFlag.Name)
	MarkFlagFilename(ctx, subCmd, flags.TemplateOutputPathFlag.Name)
	initCmd.AddCommand(subCmd)
}

func initExecFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	outputPath := flags.ValueFor[string](ctx, cmd, *flags.TemplateOutputPathFlag, false)
	template := flags.ValueFor[string](ctx, cmd, *flags.TemplateFlag, false)
	templateFilePath := flags.ValueFor[string](ctx, cmd, *flags.TemplateFilePathFlag, false)
	workspaceName := flags.ValueFor[string](ctx, cmd, *flags.TemplateWorkspaceFlag, false)

	ws := workspaceOrCurrent(ctx, workspaceName)
	if ws == nil {
		logger.Fatalf("workspace %s not found", workspaceName)
	}

	tmpl := loadFlowfileTemplate(ctx, template, templateFilePath)
	if tmpl == nil {
		logger.Fatalf("unable to load flowfile template")
	}

	flowFilename := tmpl.Name()
	if len(args) == 1 {
		flowFilename = args[0]
	}
	if err := templates.ProcessTemplate(ctx, tmpl, ws, flowFilename, outputPath); err != nil {
		logger.FatalErr(err)
	}

	logger.PlainTextSuccess(fmt.Sprintf("Template '%s' rendered successfully", flowFilename))
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

func initVaultFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	generatedKey, err := crypto.GenerateKey()
	if err != nil {
		logger.FatalErr(err)
	}
	if err = vault.RegisterEncryptionKey(generatedKey); err != nil {
		logger.FatalErr(err)
	}

	if verbosity := flags.ValueFor[int](ctx, cmd, *flags.VerbosityFlag, false); verbosity >= 0 {
		logger.PlainTextSuccess(fmt.Sprintf("Your vault encryption key is: %s", generatedKey))
		newKeyMsg := fmt.Sprintf(
			"You will need this key to modify your vault data. Store it somewhere safe!\n"+
				"Set this value to the %s environment variable if you do not want to be prompted for it every time.",
			vault.EncryptionKeyEnvVar,
		)
		logger.PlainTextInfo(newKeyMsg)
	} else {
		logger.PlainTextSuccess(fmt.Sprintf("Encryption key: %s", generatedKey))
	}
}

//nolint:lll
var initExecLong = `Add rendered executables from an executable definition template to a workspace.

The WORKSPACE_NAME is the name of the workspace to initialize the executables in.
The DEFINITION_NAME is the name of the definition to use in rendering the template.
This name will become the name of the file containing the copied executable definition.

One one of -f or -t must be provided and must point to a valid executable definition template.
The -p flag can be used to specify a sub-path within the workspace to create the executable definition and its artifacts.`
