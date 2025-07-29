package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"net/url"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/flowexec/flow/internal/cache"
)

type MCPServer struct {
	srv *mcp.Server
}

func NewMCPServer() *MCPServer {
	server := mcp.NewServer(&mcp.Implementation{Name: "flow"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "executables"}, ListExecs)
	server.AddPrompt(&mcp.Prompt{Name: "executables"}, PromptExecList)
	server.AddResource(&mcp.Resource{Name: "flowfile schema", MIMEType: "application/json", URI: flowFileSchemaURI}, ResourceHandler)

	return &MCPServer{srv: server}
}

func (s *MCPServer) Run() error {
	t := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
	if err := s.srv.Run(context.Background(), t); err != nil {
		return err
	}
	return nil
}

type ExecArgs struct {
	Workspace string `json:"workspace" jsonschema:"workspace to filter executables for"`
}

func ListExecs(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[ExecArgs]) (*mcp.CallToolResultFor[struct{}], error) {
	wsCache := cache.NewWorkspaceCache()
	execCache := cache.NewExecutableCache(wsCache)

	list, err := execCache.GetExecutableList(nil)
	if err != nil {
		return nil, err
	}
	wsArg := params.Arguments.Workspace
	if wsArg != "" {
		list.FilterByWorkspace(wsArg)
	}

	json, err := list.JSON()
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: json},
		},
	}, nil
}

func PromptExecList(_ context.Context, _ *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	wsArg := params.Arguments["name"]
	promptTxt := fmt.Sprintf("Give me a summary of all of the executables in the %s workspace", wsArg)
	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{
			{
				Role: "user", Content: &mcp.TextContent{Text: promptTxt},
			},
		},
	}, nil

}

var (
	//go:embed schemas/flowfile_schema.json
	flowFileSchema string
	//go:embed schemas/workspace_schema.json
	workspaceSchema string

	flowFileSchemaURI        = "flow://schema/flowfile"
	workspaceConfigSchemaURI = "flow://schema/workspace"
)

func ResourceHandler(_ context.Context, _ *mcp.ServerSession, params *mcp.ReadResourceParams) (*mcp.ReadResourceResult, error) {
	u, err := url.Parse(params.URI)
	if err != nil {
		return nil, err
	}

	var content string
	switch u.String() {
	case flowFileSchemaURI:
		content = flowFileSchema
	case workspaceConfigSchemaURI:
		content = workspaceSchema
	default:
		return nil, fmt.Errorf("unknown resource %s", u.String())
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{URI: params.URI, MIMEType: "application/json", Text: content},
		},
	}, nil
}
