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

	getWorkspace := mcp.NewTool("get_workspace",
		mcp.WithString("workspace_name", mcp.Required(), mcp.Description("Registered workspace name")),
		mcp.WithDescription("Get details about a registered flow workspaces"),
	)
	srv.AddTool(getWorkspace, getWorkspaceHandler)

	listExecutables := mcp.NewTool("list_executables",
		mcp.WithDescription("List executables in a workspace"),
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
		mcp.WithBoolean("sync", mcp.Description("Sync executable cache before execution")),
	)
	srv.AddTool(executeFlow, executeFlowHandler)

	getExecutable := mcp.NewTool("get_executable",
		mcp.WithDescription("Get detailed information about an executable"),
		mcp.WithString("executable_verb", mcp.Required(),
			mcp.Enum(executable.SortedValidVerbs()...),
			mcp.Description("Executable verb")),
		mcp.WithString("executable_id",
			mcp.Pattern(`^([a-zA-Z0-9_-]+(/[a-zA-Z0-9_-]+)?:)?[a-zA-Z0-9_-]+$`),
			mcp.Description("Executable ID (workspace/namespace:name or just name if using the current workspace/namespace)")),
	)
	srv.AddTool(getExecutable, getExecutableHandler)
}

func listWorkspacesHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd := exec.Command("flow", "workspace", "list", "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list workspaces: %s", output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func getWorkspaceHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	wsName, err := request.RequireString("workspace_name")
	if err != nil {
		return mcp.NewToolResultError("workspace_name is required"), nil
	}

	cmd := exec.Command("flow", "workspace", "get", wsName, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get workspaces %s: %s", wsName, output)), nil
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

func getExecutableHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	executableVerb, err := request.RequireString("executable_verb")
	if err != nil {
		return mcp.NewToolResultError("executable_verb is required"), nil
	}
	executableID := request.GetString("executable_id", "")

	cmdArgs := []string{"browse", "--output", "json", executableVerb}
	if executableID != "" {
		cmdArgs = append(cmdArgs, executableID)
	}
	cmd := exec.Command("flow", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		ref := strings.Join([]string{executableVerb, executableID}, " ")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get executable %s: %s", ref, output)), nil
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
