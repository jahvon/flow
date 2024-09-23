package internal

import (
	"fmt"

	"github.com/jahvon/tuikit/types"
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/cmd/internal/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io/executable"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/exec"
	"github.com/jahvon/flow/internal/templates"
)

func RegisterTemplateCmd(ctx *context.Context, rootCmd *cobra.Command) {
	templateCmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"tmpl", "t"},
		Short:   "Manage flowfile templates.",
	}
	registerInitTemplateCmd(ctx, templateCmd)
	registerSetTemplateCmd(ctx, templateCmd)
	registerListTemplateCmd(ctx, templateCmd)
	registerGetTemplateCmd(ctx, templateCmd)
	rootCmd.AddCommand(templateCmd)
}

func registerInitTemplateCmd(ctx *context.Context, templateCmd *cobra.Command) {
	initCmd := &cobra.Command{
		Use:     "init FLOWFILE_NAME [-w WORKSPACE ] [-o OUTPUT_DIR] [-f FILE | -t TEMPLATE]",
		Aliases: []string{"render", "run"},
		Short:   "Render a flowfile template into a workspace.",
		Long:    templateLong,
		Args:    cobra.MaximumNArgs(1),
		PreRun:  func(cmd *cobra.Command, args []string) { runner.RegisterRunner(exec.NewRunner()) },
		Run:     func(cmd *cobra.Command, args []string) { templateFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, initCmd, *flags.TemplateOutputPathFlag)
	RegisterFlag(ctx, initCmd, *flags.TemplateFlag)
	RegisterFlag(ctx, initCmd, *flags.TemplateFilePathFlag)
	RegisterFlag(ctx, initCmd, *flags.TemplateWorkspaceFlag)
	MarkFlagMutuallyExclusive(initCmd, flags.TemplateFlag.Name, flags.TemplateFilePathFlag.Name)
	MarkOneFlagRequired(initCmd, flags.TemplateFlag.Name, flags.TemplateFilePathFlag.Name)
	MarkFlagFilename(ctx, initCmd, flags.TemplateFilePathFlag.Name)
	MarkFlagFilename(ctx, initCmd, flags.TemplateOutputPathFlag.Name)
	templateCmd.AddCommand(initCmd)
}

func templateFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
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

func registerSetTemplateCmd(ctx *context.Context, templateCmd *cobra.Command) {
	setCmd := &cobra.Command{
		Use:     "set NAME DEFINITION_TEMPLATE_PATH",
		Aliases: []string{"new"},
		Short:   "Register a flowfile template.",
		Args:    cobra.ExactArgs(2),
		Run:     func(cmd *cobra.Command, args []string) { setTemplateFunc(ctx, cmd, args) },
	}
	templateCmd.AddCommand(setCmd)
}

func setTemplateFunc(ctx *context.Context, _ *cobra.Command, args []string) {
	logger := ctx.Logger
	name := args[0]
	flowFilePath := args[1]
	loadedTemplates, err := filesystem.LoadFlowFileTemplate(name, flowFilePath)
	if err != nil {
		logger.FatalErr(err)
	}
	if err := loadedTemplates.Validate(); err != nil {
		logger.FatalErr(err)
	}
	userConfig := ctx.Config
	if userConfig.Templates == nil {
		userConfig.Templates = map[string]string{}
	}
	userConfig.Templates[name] = flowFilePath
	if err := filesystem.WriteConfig(userConfig); err != nil {
		logger.FatalErr(err)
	}
	logger.PlainTextSuccess(fmt.Sprintf("Template %s set to %s", name, flowFilePath))
}

func registerListTemplateCmd(ctx *context.Context, templateCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "View a list of registered flowfile templates.",
		Args:    cobra.NoArgs,
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
		Run:     func(cmd *cobra.Command, args []string) { listTemplateFunc(ctx, cmd, args) },
	}
	RegisterFlag(ctx, listCmd, *flags.OutputFormatFlag)
	templateCmd.AddCommand(listCmd)
}

func listTemplateFunc(ctx *context.Context, cmd *cobra.Command, _ []string) {
	logger := ctx.Logger
	// TODO: include unregistered templates within the current ws
	tmpls, err := filesystem.LoadFlowFileTemplates(ctx.Config.Templates)
	if err != nil {
		logger.FatalErr(err)
	}

	outputFormat := flags.ValueFor[string](ctx, cmd, *flags.OutputFormatFlag, false)
	if TUIEnabled(ctx, cmd) {
		view := executable.NewTemplateListView(
			ctx, tmpls, types.Format(outputFormat),
			func(name string) error {
				tmpl := tmpls.Find(name)
				if tmpl == nil {
					return fmt.Errorf("template %s not found", name)
				}
				ws := ctx.CurrentWorkspace
				// TODO: support specifying a path/name
				if err := templates.ProcessTemplate(ctx, tmpl, ws, tmpl.Name(), "//"); err != nil {
					return err
				}
				logger.PlainTextSuccess("Template rendered successfully")
				return nil
			},
		)
		SetView(ctx, cmd, view)
	} else {
		executable.PrintTemplateList(logger, outputFormat, tmpls)
	}
}

func registerGetTemplateCmd(ctx *context.Context, getCmd *cobra.Command) {
	templateCmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"find"},
		Short:   "View a flowfile template's documentation. Either it's registered name or file path can be used.",
		PreRun:  func(cmd *cobra.Command, args []string) { StartTUI(ctx, cmd) },
		PostRun: func(cmd *cobra.Command, args []string) { WaitForTUI(ctx, cmd) },
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
	if TUIEnabled(ctx, cmd) {
		runFunc := func(ref string) error { return runByRef(ctx, cmd, ref) }
		view := executable.NewTemplateView(ctx, tmpl, types.Format(outputFormat), runFunc)
		SetView(ctx, cmd, view)
	} else {
		executable.PrintTemplate(logger, outputFormat, tmpl)
	}
}

var templateLong = `Add rendered executables from a flowfile template to a workspace.

The WORKSPACE_NAME is the name of the workspace to initialize the flowfile template in.
The FLOWFILE_NAME is the name to give the flowfile (if applicable) when rendering its template.

One one of -f or -t must be provided and must point to a valid flowfile template.
The -o flag can be used to specify an output path within the workspace to create the flowfile and its artifacts in.`
