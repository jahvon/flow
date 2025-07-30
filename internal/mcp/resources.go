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
	//go:embed resources/concepts-guide.md
	conceptsMD string

	//go:embed resources/file-types-guide.md
	fileTypesMD string

	//go:embed resources/flowfile_schema.json
	flowFileSchema string

	//go:embed resources/template_schema.json
	templateSchema string

	//go:embed resources/workspace_schema.json
	workspaceSchema string
)

func addServerResources(srv *server.MCPServer) {
	getCtx := mcp.NewResource(
		"flow://context/current",
		"context",
		mcp.WithResourceDescription("Current flow execution context (workspace, namespace, vault)"),
		mcp.WithMIMEType("application/json"))
	srv.AddResource(getCtx, getContextResourceHandler)

	getWorkspace := mcp.NewResourceTemplate(
		"flow://workspace/{name}",
		"workspace",
		mcp.WithTemplateDescription("Flow workspace configuration and details"))
	srv.AddResourceTemplate(getWorkspace, getWorkspaceResourceHandler)

	getWorkspaceExecutables := mcp.NewResourceTemplate(
		"flow://workspace/{name}/executables",
		"workspace_executables",
		mcp.WithTemplateDescription("List of Flow executables for a given workspace"))
	srv.AddResourceTemplate(getWorkspaceExecutables, getWorkspaceExecutablesHandler)

	getExecutable := mcp.NewResourceTemplate(
		"flow://executable/{ref}",
		"executable",
		mcp.WithTemplateDescription("Flow executable details by reference"))
	srv.AddResourceTemplate(getExecutable, getExecutableResourceHandler)

	getFlowConcepts := mcp.NewResource(
		"flow://guide/concepts",
		"flow_concepts",
		mcp.WithResourceDescription("Information about flow and it's usage"),
		mcp.WithMIMEType("text/markdown"))
	srv.AddResource(getFlowConcepts, getFlowConceptsResourceHandler)

	getFlowFileTypes := mcp.NewResource(
		"flow://guide/file-types",
		"flow_file_types",
		mcp.WithResourceDescription("Information about flow file types and their usage"),
		mcp.WithMIMEType("text/markdown"))
	srv.AddResource(getFlowFileTypes, getFlowFileTypesResourceHandler)

	getFlowFileSchema := mcp.NewResource(
		"flow://schema/flowfile",
		"flowfile_schema",
		mcp.WithResourceDescription("Flow file (*.flow, *.flow.yaml, *.flow.yml) schema"),
		mcp.WithMIMEType("application/json"))
	srv.AddResource(getFlowFileSchema, getFlowFileSchemaResourceHandler)

	getTemplateSchema := mcp.NewResource(
		"flow://schema/template",
		"template_schema",
		mcp.WithResourceDescription("Flow template (*.flow.tmpl, *.flow.yaml.tmpl, *.flow.yml.tmpl) schema"),
		mcp.WithMIMEType("application/json"))
	srv.AddResource(getTemplateSchema, getTemplateSchemaResourceHandler)

	getWorkspaceSchema := mcp.NewResource(
		"flow://schema/workspace",
		"workspace_schema",
		mcp.WithResourceDescription("Flow workspace configuration (flow.yaml) schema"),
		mcp.WithMIMEType("application/json"))
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

func getExecutableResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	executableVerb, executableID, err := extractExecutableRef(request.Params.URI)
	if err != nil {
		return nil, err
	}

	cmdArgs := []string{"browse", "--output", "json", executableVerb}
	if executableID != "" {
		cmdArgs = append(cmdArgs, executableID)
	}
	cmd := exec.Command("flow", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		ref := strings.Join([]string{executableVerb, executableID}, " ")
		return nil, fmt.Errorf("%s executable details retrieval failed: %s", ref, output)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(output),
		},
	}, nil
}

func getFlowConceptsResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/markdown",
			Text:     conceptsMD,
		},
	}, nil
}

func getFlowFileTypesResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/markdown",
			Text:     fileTypesMD,
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

func getTemplateSchemaResourceHandler(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     templateSchema,
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

func extractExecutableRef(uri string) (string, string, error) {
	if len(uri) < len("flow://executable/") {
		return "", "", fmt.Errorf("invalid executable URI: %s", uri)
	}

	ref := uri[len("flow://executable/"):]
	if ref == "" {
		return "", "", fmt.Errorf("executable reference cannot be empty in URI: %s", uri)
	}

	parts := strings.SplitN(ref, " ", 2)
	switch len(parts) {
	case 1:
		return parts[0], "", nil
	case 2:
		return parts[0], parts[1], nil
	default:
		return "", "", fmt.Errorf("invalid executable reference format in URI: %s", uri)
	}
}
