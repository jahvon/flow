package mcp

import (
	_ "embed"

	"github.com/mark3labs/mcp-go/server"

	"github.com/flowexec/flow/internal/context"
)

type Server struct {
	ctx *context.Context
	srv *server.MCPServer
}

func NewServer(ctx *context.Context) *Server {
	srv := server.NewMCPServer(
		"Flow",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithPromptCapabilities(false),
	)
	addServerTools(srv)
	addServerPrompts(srv)

	return &Server{ctx: ctx, srv: srv}
}

func (s *Server) Run() error {
	return server.ServeStdio(s.srv)
}
