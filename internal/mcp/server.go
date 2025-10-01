package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"github.com/joao-costa/multichatmcp/internal/messenger"
)

// Server wraps the MCP server with messenger functionality
type Server struct {
	mcpServer *server.MCPServer
	messenger messenger.Messenger
}

// NewServer creates a new MCP server with messenger-specific tools
func NewServer(m messenger.Messenger) *Server {
	s := &Server{
		messenger: m,
	}

	// Create MCP server with dynamic name based on messenger type
	serverName := fmt.Sprintf("MultiChat MCP Server (%s)", m.GetMessengerName())
	mcpServer := server.NewMCPServer(
		serverName,
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Let the messenger register its own tools
	m.RegisterMCPTools(mcpServer)

	s.mcpServer = mcpServer
	return s
}

// Serve starts the MCP server
func (s *Server) Serve() error {
	return server.ServeStdio(s.mcpServer)
}
