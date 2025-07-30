package mcp

import (
	_ "embed"

	"github.com/mark3labs/mcp-go/server"

	"github.com/flowexec/flow/internal/context"
)

type MCPServer struct {
	ctx *context.Context
	srv *server.MCPServer
}

func NewMCPServer(ctx *context.Context) *MCPServer {
	srv := server.NewMCPServer(
		"Flow",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithResourceCapabilities(false, true),
		server.WithPromptCapabilities(false),
	)

	addServerResources(srv)
	addServerTools(srv)
	addServerPrompts(srv)

	return &MCPServer{ctx: ctx, srv: srv}
}

func (s *MCPServer) Run() error {
	return server.ServeStdio(s.srv)
}
