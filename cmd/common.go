package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/ui"
	"github.com/jahvon/flow/internal/io/ui/views"
)

var (
	curCtx   *context.Context
	termView *views.StandardTermView
)

func interactiveUIEnabled() bool {
	disabled := getPersistentFlagValue[bool](rootCmd, *flags.NonInteractiveFlag)
	return !disabled && curCtx.UserConfig.Interactive != nil && curCtx.UserConfig.Interactive.Enabled
}

func handleError(err error) {
	if interactiveUIEnabled() {
		curCtx.App.HandleInternalError(err)
		for {
			if curCtx.App.Ready() {
				curCtx.App.Finalize()
				time.Sleep(1 * time.Second)
				curCtx.CancelFunc()
				os.Exit(1)
			}
		}
	} else {
		curCtx.Logger.FatalErr(err)
	}
}

func processUserInput(inputs ...*views.TextInput) map[string]string {
	collectedVals := make(map[string]string)
	if interactiveUIEnabled() {
		if termView == nil {
			handleError(fmt.Errorf("unable to process user input"))
		}
		termView.StartProcessingUserInputs(inputs...)
		if err := termView.WaitForTextInputs(); err != nil {
			handleError(err)
		}
		for _, input := range termView.GetTextInputs() {
			collectedVals[input.Key] = input.Value()
		}
	} else {
		for _, input := range inputs {
			val := io.ProcessUserInput(input.Prompt)
			collectedVals[input.Key] = val
		}
	}
	return collectedVals
}

func processUserConfirmation(prompt string) bool {
	if interactiveUIEnabled() {
		termView.StartProcessingUserInputs(&views.TextInput{
			Key:    "confirmation",
			Prompt: prompt + " (y/n) ",
		})
		if err := termView.WaitForTextInputs(); err != nil {
			handleError(err)
		}
		response := termView.GetTextInputs()[0].Value()
		return response == "y" || response == "Y"
	}
	response := io.ProcessUserInput(prompt + " (y/n) ")
	return response == "y" || response == "Y"
}

func setTermView(cmd *cobra.Command, args []string) {
	startApp(cmd, args)
	if interactiveUIEnabled() {
		go func() {
			for {
				if curCtx.App.Ready() {
					var ok bool
					termView, ok = (views.NewTermView(curCtx.App, curCtx.Logger)).(*views.StandardTermView)
					if !ok {
						handleError(fmt.Errorf("unable to set term view"))
					}
					curCtx.App.BuildAndSetView(termView)
					break
				}
			}
		}()
	}
}

func startApp(_ *cobra.Command, _ []string) {
	enabled := interactiveUIEnabled()
	if enabled && curCtx.App == nil {
		app := ui.StartApplication(curCtx.Ctx, curCtx.CancelFunc)
		app.SetContext(curCtx.UserConfig.CurrentWorkspace, curCtx.UserConfig.CurrentNamespace)
		curCtx.App = app
	} else if !enabled {
		curCtx.Logger.SetBackground(false)
	}
}

func waitForExit(_ *cobra.Command, _ []string) {
	if interactiveUIEnabled() && curCtx.App != nil {
		timeout := time.After(30 * time.Minute)
		select {
		case <-curCtx.Ctx.Done():
			return
		case <-timeout:
			panic("interactive wait timeout")
		}
	}
}

func exitApp(_ *cobra.Command, _ []string) {
	if interactiveUIEnabled() {
		for {
			if !curCtx.Logger.PendingRead() {
				curCtx.App.Finalize()
				time.Sleep(1 * time.Second)
				curCtx.CancelFunc()
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func GenerateMarkdownTree(dir string) error {
	return doc.GenMarkdownTree(rootCmd, dir)
}
