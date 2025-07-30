package mcp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/flowexec/flow/types/executable"
)

func addServerTools(srv *server.MCPServer) {
	listWorkspaces := mcp.NewTool("list_workspaces",
		mcp.WithDescription("List all registered flow workspaces"),
	)
	srv.AddTool(listWorkspaces, listWorkspacesHandler)

	listExecutables := mcp.NewTool("list_executables",
		mcp.WithDescription("List and filter executables across all workspaces"),
		mcp.WithString("workspace", mcp.Description("Workspace name (optional)")),
		mcp.WithString("namespace", mcp.Description("Namespace filter (optional)")),
		mcp.WithString("verb", mcp.Description("Verb filter (optional)")),
		mcp.WithString("keyword", mcp.Description("Keyword filter (optional)")),
		mcp.WithString("tag", mcp.Description("Tag filter (optional)")),
	)
	srv.AddTool(listExecutables, listExecutablesHandler)

	executeFlow := mcp.NewTool("execute_flow",
		mcp.WithDescription("Execute a flow executable"),
		mcp.WithString("executable_verb", mcp.Required(),
			mcp.Enum(executable.SortedValidVerbs()...),
			mcp.Description("Executable verb")),
		mcp.WithString("executable_id",
			mcp.Pattern(`^([a-zA-Z0-9_-]+(/[a-zA-Z0-9_-]+)?:)?[a-zA-Z0-9_-]+$`),
			mcp.Description("Executable ID (workspace/namespace:name or just name if using the current workspace/namespace)")),
		mcp.WithString("args", mcp.Description("Arguments to pass")),
		mcp.WithBoolean("sync", mcp.Description("Sync executable changes before execution")),
	)
	srv.AddTool(executeFlow, executeFlowHandler)

	getExecutionLogs := mcp.NewTool("get_execution_logs",
		mcp.WithDescription("Get a list of the recent flow execution logs"),
		mcp.WithBoolean("last", mcp.Description("Get only the last execution logs")))
	srv.AddTool(getExecutionLogs, getExecutionLogsHandler)

	sync := mcp.NewTool("sync_executables",
		mcp.WithDescription("Sync the flow workspace and executable state"))
	srv.AddTool(sync, syncStateHandler)
}

func listWorkspacesHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd := exec.Command("flow", "workspace", "list", "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list workspaces: %s", output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func listExecutablesHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	wsFilter := request.GetString("workspace", executable.WildcardWorkspace)
	nsFilter := request.GetString("namespace", executable.WildcardNamespace)
	verbFilter := request.GetString("verb", "")
	keywordFilter := request.GetString("keyword", "")
	tagFilter := request.GetString("tag", "")

	cmdArgs := []string{"browse", "--output", "json", "--workspace", wsFilter, "--namespace", nsFilter}
	if verbFilter != "" {
		cmdArgs = append(cmdArgs, "--verb", verbFilter)
	}
	if keywordFilter != "" {
		cmdArgs = append(cmdArgs, "--filter", keywordFilter)
	}
	if tagFilter != "" {
		cmdArgs = append(cmdArgs, "--tag", tagFilter)
	}

	cmd := exec.Command("flow", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list executables: %s", output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func executeFlowHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	executableVerb, err := request.RequireString("executable_verb")
	if err != nil {
		return mcp.NewToolResultError("executable_verb is required"), nil
	}
	executableID := request.GetString("executable_id", "")

	args := request.GetString("args", "")
	sync := request.GetBool("sync", false)

	cmdArgs := []string{executableVerb}
	if executableID != "" {
		cmdArgs = append(cmdArgs, executableID)
	}
	if args != "" {
		cmdArgs = append(cmdArgs, strings.Fields(args)...)
	}
	if sync {
		cmdArgs = append(cmdArgs, "--sync")
	}

	cmd := exec.Command("flow", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		ref := strings.Join([]string{executableVerb, executableID}, " ")
		return mcp.NewToolResultError(fmt.Sprintf("%s execution failed: %s", ref, string(output))), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func getExecutionLogsHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	last := request.GetBool("last", false)
	cmdArgs := []string{"logs", "--output", "json"}
	if last {
		cmdArgs = append(cmdArgs, "--last")
	}
	cmd := exec.Command("flow", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get flow execution logs: %s", output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func syncStateHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd := exec.Command("flow", "sync")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to sync flow's state: %s", output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}
