package mcp

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"

	"github.com/flowexec/flow/internal/filesystem"
)

var (
	//go:embed resources/info.md
	infoMD string

	//go:embed resources/flowfile_schema.json
	flowFileSchema string

	//go:embed resources/workspace_schema.json
	workspaceSchema string

	flowFileSchemaURI        = "flow://schema/flowfile"
	workspaceConfigSchemaURI = "flow://schema/workspace"
)

func addServerResources(srv *server.MCPServer) {
	getCtx := mcp.NewResource(
		"flow://context/current",
		"context",
		mcp.WithResourceDescription("Current flow execution context (workspace, namespace, vault)"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(getCtx, getContextResourceHandler)

	getWorkspace := mcp.NewResourceTemplate(
		"flow://workspace/{name}",
		"workspace",
		mcp.WithTemplateDescription("Flow workspace configuration and details"),
	)
	srv.AddResourceTemplate(getWorkspace, getWorkspaceResourceHandler)

	getWorkspaceExecutables := mcp.NewResourceTemplate(
		"flow://workspace/{name}/executables",
		"workspace_executables",
		mcp.WithTemplateDescription("Flow executables for a given workspace"),
	)
	srv.AddResourceTemplate(getWorkspaceExecutables, getWorkspaceExecutablesHandler)

	getFlowInfo := mcp.NewResource(
		"flow://info",
		"flow_info",
		mcp.WithResourceDescription("Information about Flow and it's usage"),
		mcp.WithMIMEType("text/markdown"),
	)
	srv.AddResource(getFlowInfo, getFlowInfoResourceHandler)

	getFlowFileSchema := mcp.NewResource(
		flowFileSchemaURI,
		"flowfile_schema",
		mcp.WithResourceDescription("Flow file (*.flow, *.flow.yaml, *.flow.yml) schema"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(getFlowFileSchema, getFlowFileSchemaResourceHandler)

	getWorkspaceSchema := mcp.NewResource(
		workspaceConfigSchemaURI,
		"workspace_schema",
		mcp.WithResourceDescription("Flow workspace configuration (flow.yaml) schema"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(getWorkspaceSchema, getWorkspaceSchemaResourceHandler)
}

func getContextResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	cfg, err := filesystem.LoadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load user config")
	}
	cfg.SetDefaults()

	wsName, err := cfg.CurrentWorkspaceName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current workspace name")
	}

	contextData := map[string]interface{}{
		"workspace":     wsName,
		"namespace":     cfg.CurrentNamespace,
		"vault":         cfg.CurrentVault,
		"workspaceMode": cfg.WorkspaceMode,
		"workspacePath": cfg.Workspaces[cfg.CurrentWorkspace],
	}
	jsonData, err := json.Marshal(contextData)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func getWorkspaceResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	wsName, err := extractWsName(request.Params.URI)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("flow", "workspace", "get", wsName, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace %s: %s", wsName, output)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(output),
		},
	}, nil
}

func getWorkspaceExecutablesHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	wsName, err := extractWsName(request.Params.URI)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("flow", "browse", "--output", "json", "--workspace", wsName, "--all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list executables for workspace %s: %s", wsName, output)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(output),
		},
	}, nil
}

func extractWsName(uri string) (string, error) {
	// Assuming the URI is in the format "flow://workspace/{name}" or "flow://workspace/{name}/executables"
	if len(uri) < len("flow://workspace/") {
		return "", fmt.Errorf("invalid workspace URI: %s", uri)
	}

	wsName := uri[len("flow://workspace/"):]
	strings.TrimSuffix(wsName, "/executables")
	if wsName == "" {
		return "", fmt.Errorf("workspace name cannot be empty in URI: %s", uri)
	}
	return wsName, nil
}

func getFlowInfoResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/markdown",
			Text:     infoMD,
		},
	}, nil
}

func getFlowFileSchemaResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     flowFileSchema,
		},
	}, nil
}

func getWorkspaceSchemaResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     workspaceSchema,
		},
	}, nil
}
