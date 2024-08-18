package internal

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/jahvon/tuikit/components"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/cmd/internal/interactive"
	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	configio "github.com/jahvon/flow/internal/io/config"
	executableio "github.com/jahvon/flow/internal/io/executable"
	workspaceio "github.com/jahvon/flow/internal/io/workspace"
	"github.com/jahvon/flow/internal/vault"
	"github.com/jahvon/flow/types/executable"
)

func RegisterGetCmd(ctx *context.Context, rootCmd *cobra.Command) {
	getCmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Print a flow entity.",
	}
	registerGetConfigCmd(ctx, getCmd)
	registerGetWsCmd(ctx, getCmd)
	registerGetExecCmd(ctx, getCmd)
	registerGetSecretCmd(ctx, getCmd)
	registerGetTemplateCmd(ctx, getCmd)
	rootCmd.AddCommand(getCmd)
}

func registerGetConfigCmd(ctx *context.Context, getCmd *cobra.Command) {
	configCmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Print the current global configuration values.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveContainer(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { interactive.WaitForExit(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { getConfigFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, configCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(configCmd)
}

func getConfigFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	userConfig := ctx.Config
	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if interactive.UIEnabled(ctx, cmd) {
		view := configio.NewUserConfigView(ctx.InteractiveContainer, *userConfig, components.Format(outputFormat))
		ctx.InteractiveContainer.SetView(view)
	} else {
		configio.PrintUserConfig(logger, outputFormat, userConfig)
	}
}

func registerGetWsCmd(ctx *context.Context, getCmd *cobra.Command) {
	wsCmd := &cobra.Command{
		Use:     "workspace [NAME]",
		Aliases: []string{"ws"},
		Short:   "Print a workspace's configuration. If the name is omitted, the current workspace is used.",
		Args:    cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return maps.Keys(ctx.Config.Workspaces), cobra.ShellCompDirectiveNoFileComp
		},
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveContainer(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { interactive.WaitForExit(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { getWsFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, wsCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(wsCmd)
}

func getWsFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	var workspaceName, wsPath string
	if len(args) == 1 {
		workspaceName = args[0]
		wsPath = ctx.Config.Workspaces[workspaceName]
	} else {
		workspaceName = ctx.CurrentWorkspace.AssignedName()
		wsPath = ctx.CurrentWorkspace.Location()
	}

	wsCfg, err := filesystem.LoadWorkspaceConfig(workspaceName, wsPath)
	if err != nil {
		logger.FatalErr(errors.Wrap(err, "failure loading workspace config"))
	} else if wsCfg == nil {
		logger.Fatalf("config not found for workspace %s", workspaceName)
	}

	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if interactive.UIEnabled(ctx, cmd) {
		view := workspaceio.NewWorkspaceView(ctx, wsCfg, components.Format(outputFormat))
		ctx.InteractiveContainer.SetView(view)
	} else {
		workspaceio.PrintWorkspaceConfig(logger, outputFormat, wsCfg)
	}
}

func registerGetExecCmd(ctx *context.Context, getCmd *cobra.Command) {
	execCmd := &cobra.Command{
		Use:     "executable VERB ID",
		Aliases: []string{"exec"},
		Short:   "Print an executable flow by reference.",
		Long: "Print an executable by the executable's verb and ID.\nThe target executable's ID should be in the  " +
			"form of 'ws/ns:name' and the verb should match the target executable's verb or one of its aliases.\n\n" +
			fmt.Sprintf(
				"See %s for more information on executable verbs.\n"+
					"See %s for more information on executable IDs.",
				io.TypesDocsURL("flowfile", "ExecutableVerb"),
				io.TypesDocsURL("flowfile", "ExecutableRef"),
			),
		Args:    cobra.ExactArgs(2),
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveContainer(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { interactive.WaitForExit(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { getExecFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, execCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(execCmd)
}

func getExecFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	verbStr := args[0]
	verb := executable.Verb(verbStr)
	if err := verb.Validate(); err != nil {
		logger.FatalErr(err)
	}
	id := args[1]
	ws, ns, name := executable.ParseExecutableID(id)
	if ws == "" {
		ws = ctx.CurrentWorkspace.AssignedName()
	}
	if ns == "" && ctx.Config.CurrentNamespace != "" {
		ns = ctx.Config.CurrentNamespace
	}
	id = executable.NewExecutableID(ws, ns, name)
	ref := executable.NewRef(id, verb)

	exec, err := ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	if err != nil && errors.Is(cache.NewExecutableNotFoundError(ref.String()), err) {
		logger.Debugf("Executable %s not found in cache, syncing cache", ref)
		if err := ctx.ExecutableCache.Update(logger); err != nil {
			logger.FatalErr(err)
		}
		exec, err = ctx.ExecutableCache.GetExecutableByRef(logger, ref)
	}
	if err != nil {
		logger.FatalErr(err)
	} else if exec == nil {
		logger.Fatalf("executable %s not found", ref)
	}

	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if interactive.UIEnabled(ctx, cmd) {
		runFunc := func(ref string) error { return runByRef(ctx, cmd, ref) }
		view := executableio.NewExecutableView(ctx, *exec, components.Format(outputFormat), runFunc)
		ctx.InteractiveContainer.SetView(view)
	} else {
		executableio.PrintExecutable(logger, outputFormat, exec)
	}
}

func registerGetSecretCmd(ctx *context.Context, getCmd *cobra.Command) {
	secretCmd := &cobra.Command{
		Use:     "secret REFERENCE",
		Aliases: []string{"scrt"},
		Short:   "Print the value of a secret in the flow secret vault.",
		Args:    cobra.ExactArgs(1),
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveCommand(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { getSecretFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, secretCmd, *flags.OutputSecretAsPlainTextFlag)
	RegisterFlag(ctx, secretCmd, *flags.CopyFlag)
	getCmd.AddCommand(secretCmd)
}

func getSecretFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	logger := ctx.Logger
	reference := args[0]
	asPlainText := flags.ValueFor[bool](ctx, cmd, *flags.OutputSecretAsPlainTextFlag, false)
	copyValue := flags.ValueFor[bool](ctx, cmd, *flags.CopyFlag, false)

	v := vault.NewVault(logger)
	secret, err := v.GetSecret(reference)
	if err != nil {
		logger.FatalErr(err)
	}

	if asPlainText {
		logger.PlainTextInfo(secret.PlainTextString())
	} else {
		logger.PlainTextInfo(secret.String())
	}

	if copyValue {
		if err := clipboard.WriteAll(secret.PlainTextString()); err != nil {
			logger.Error(err, "\nunable to copy secret value to clipboard")
		} else {
			logger.PlainTextSuccess("\ncopied secret value to clipboard")
		}
	}
}

func registerGetTemplateCmd(ctx *context.Context, getCmd *cobra.Command) {
	templateCmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"tmpl"},
		Short:   "Print a flowfile template using it's registered name or file path.",
		PreRun:  func(cmd *cobra.Command, args []string) { interactive.InitInteractiveContainer(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { interactive.WaitForExit(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { getTemplateFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, templateCmd, *flags.TemplateFlag)
	RegisterFlag(ctx, templateCmd, *flags.TemplateFilePathFlag)
	MarkOneFlagRequired(templateCmd, flags.TemplateFlag.Name, flags.TemplateFilePathFlag.Name)
	RegisterFlag(ctx, templateCmd, *flags.OutputFormatFlag)
	getCmd.AddCommand(templateCmd)
}

func getTemplateFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	template := flags.ValueFor[string](ctx, cmd, *flags.TemplateFlag, false)
	templateFilePath := flags.ValueFor[string](ctx, cmd, *flags.TemplateFilePathFlag, false)

	tmpl := loadFlowfileTemplate(ctx, template, templateFilePath)
	if tmpl == nil {
		logger.Fatalf("unable to load flowfile template")
	}

	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if interactive.UIEnabled(ctx, cmd) {
		runFunc := func(ref string) error { return runByRef(ctx, cmd, ref) }
		view := executableio.NewTemplateView(ctx, tmpl, components.Format(outputFormat), runFunc)
		ctx.InteractiveContainer.SetView(view)
	} else {
		executableio.PrintTemplate(logger, outputFormat, tmpl)
	}
}
