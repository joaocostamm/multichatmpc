package messenger

import (
	"context"

	"github.com/mark3labs/mcp-go/server"
)

// Messenger defines the minimal interface that all messenger implementations must follow
type Messenger interface {
	// Connect establishes connection to the messenger service
	Connect(ctx context.Context) error

	// Disconnect closes the connection
	Disconnect() error

	// IsConnected returns the connection status
	IsConnected() bool

	// RegisterMCPTools registers messenger-specific MCP tools with the server
	// Each messenger implementation defines its own set of operations
	RegisterMCPTools(mcpServer *server.MCPServer)

	// GetMessengerName returns the name of the messenger platform (e.g., "whatsapp", "teams")
	GetMessengerName() string
}
