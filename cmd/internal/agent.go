package internal

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/context"
	agent2 "github.com/jahvon/flow/internal/services/agent"
)

func RegisterAgentCmd(ctx *context.Context, rootCmd *cobra.Command) {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage the flow background agent",
		Long:  "Manage the flow background agent for running executables in the background or on a schedule",
	}

	registerAgentInstallCmd(ctx, agentCmd)
	registerAgentUninstallCmd(ctx, agentCmd)
	registerAgentStartCmd(ctx, agentCmd)
	registerAgentStopCmd(ctx, agentCmd)
	registerAgentStatusCmd(ctx, agentCmd)
	registerAgentScheduleCmd(ctx, agentCmd)
	registerAgentListScheduledCmd(ctx, agentCmd)
	registerAgentUnscheduleCmd(ctx, agentCmd)

	rootCmd.AddCommand(agentCmd)
}

func registerAgentInstallCmd(ctx *context.Context, agentCmd *cobra.Command) {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install the flow agent as a system service",
		Run:   func(cmd *cobra.Command, args []string) { agentInstallFunc(ctx, cmd, args) },
	}
	agentCmd.AddCommand(installCmd)
}

func agentInstallFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement agent installation command
	agent, err := agent2.NewAgent(ctx.Logger, ctx.ExecutableCache, nil)
	if err != nil {
		panic(err)
	}
	agent.Run()
}

func registerAgentUninstallCmd(ctx *context.Context, agentCmd *cobra.Command) {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall the flow agent",
		Run:   func(cmd *cobra.Command, args []string) { agentUninstallFunc(ctx, cmd, args) },
	}
	agentCmd.AddCommand(uninstallCmd)
}

func agentUninstallFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement agent uninstallation command
}

func registerAgentStartCmd(ctx *context.Context, agentCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the flow agent",
		Run:   func(cmd *cobra.Command, args []string) { agentStartFunc(ctx, cmd, args) },
	}
	agentCmd.AddCommand(startCmd)
}

func agentStartFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement agent start command
}

func registerAgentStopCmd(ctx *context.Context, agentCmd *cobra.Command) {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the flow agent",
		Run:   func(cmd *cobra.Command, args []string) { agentStopFunc(ctx, cmd, args) },
	}
	agentCmd.AddCommand(stopCmd)
}

func agentStopFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement agent stop command
}

func registerAgentStatusCmd(ctx *context.Context, agentCmd *cobra.Command) {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check the status of the flow agent",
		Run:   func(cmd *cobra.Command, args []string) { agentStatusFunc(ctx, cmd, args) },
	}
	agentCmd.AddCommand(statusCmd)
}

func agentStatusFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement agent status command
}

func registerAgentScheduleCmd(ctx *context.Context, agentCmd *cobra.Command) {
	scheduleCmd := &cobra.Command{
		Use:   "schedule EXECUTABLE_ID CRON_EXPRESSION",
		Short: "Schedule an executable to run on a cron schedule",
		Args:  cobra.ExactArgs(2),
		Run:   func(cmd *cobra.Command, args []string) { agentScheduleFunc(ctx, cmd, args) },
	}
	// Add flags for env vars and arguments to be passed to the executable
	agentCmd.AddCommand(scheduleCmd)
}

func agentScheduleFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement schedule command
}

func registerAgentListScheduledCmd(ctx *context.Context, agentCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all scheduled executables",
		Run:   func(cmd *cobra.Command, args []string) { agentListFunc(ctx, cmd, args) },
	}
	agentCmd.AddCommand(listCmd)
}

func agentListFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement list command
}

func registerAgentUnscheduleCmd(ctx *context.Context, agentCmd *cobra.Command) {
	unscheduleCmd := &cobra.Command{
		Use:   "unschedule TASK_ID",
		Short: "Unschedule a previously scheduled executable",
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { agentUnscheduleFunc(ctx, cmd, args) },
	}
	agentCmd.AddCommand(unscheduleCmd)
}

func agentUnscheduleFunc(ctx *context.Context, cmd *cobra.Command, args []string) {
	// Implement unschedule command
}
