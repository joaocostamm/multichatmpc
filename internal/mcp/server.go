package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/joao-costa/multichatmcp/internal/messenger"
)

// Server wraps the MCP server with messenger functionality
type Server struct {
	mcpServer *server.MCPServer
	messenger messenger.Messenger
}

// NewServer creates a new MCP server
func NewServer(m messenger.Messenger) *Server {
	s := &Server{
		messenger: m,
	}

	mcpServer := server.NewMCPServer(
		"MultiChat MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register all tools
	s.registerTools(mcpServer)

	s.mcpServer = mcpServer
	return s
}

// registerTools registers all MCP tools
func (s *Server) registerTools(mcpServer *server.MCPServer) {
	// search_contacts
	mcpServer.AddTool(mcp.Tool{
		Name:        "search_contacts",
		Description: "Search for contacts by name or phone number",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search term to match against contact names or phone numbers",
				},
			},
			Required: []string{"query"},
		},
	}, s.handleSearchContacts)

	// list_messages
	mcpServer.AddTool(mcp.Tool{
		Name:        "list_messages",
		Description: "Retrieve messages with optional filters (e.g. time, sender) and context",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"after": map[string]interface{}{
					"type":        "string",
					"description": "ISO-8601 formatted date to only return messages after this date",
				},
				"before": map[string]interface{}{
					"type":        "string",
					"description": "ISO-8601 formatted date to only return messages before this date",
				},
				"sender_jid": map[string]interface{}{
					"type":        "string",
					"description": "Filter messages by sender JID",
				},
				"chat_jid": map[string]interface{}{
					"type":        "string",
					"description": "Filter messages by chat JID",
				},
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search term to filter messages by content",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of messages to return",
					"default":     20,
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "Page number for pagination",
					"default":     0,
				},
			},
		},
	}, s.handleListMessages)

	// list_chats
	mcpServer.AddTool(mcp.Tool{
		Name:        "list_chats",
		Description: "List available chats with metadata (name, JID, last message)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of chats to return",
					"default":     20,
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "Page number for pagination",
					"default":     0,
				},
			},
		},
	}, s.handleListChats)

	// get_chat
	mcpServer.AddTool(mcp.Tool{
		Name:        "get_chat",
		Description: "Get information about a specific chat (metadata, messages)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"chat_jid": map[string]interface{}{
					"type":        "string",
					"description": "The JID of the chat to retrieve",
				},
			},
			Required: []string{"chat_jid"},
		},
	}, s.handleGetChat)

	// get_direct_chat_by_contact
	mcpServer.AddTool(mcp.Tool{
		Name:        "get_direct_chat_by_contact",
		Description: "Find a direct chat with a specific contact by phone number",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"phone_number": map[string]interface{}{
					"type":        "string",
					"description": "Phone number of the contact (with country code, no + or spaces)",
				},
			},
			Required: []string{"phone_number"},
		},
	}, s.handleGetDirectChatByContact)

	// get_contact_chats
	mcpServer.AddTool(mcp.Tool{
		Name:        "get_contact_chats",
		Description: "List all chats involving a specific contact",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"contact_jid": map[string]interface{}{
					"type":        "string",
					"description": "The JID of the contact",
				},
			},
			Required: []string{"contact_jid"},
		},
	}, s.handleGetContactChats)

	// send_message
	mcpServer.AddTool(mcp.Tool{
		Name:        "send_message",
		Description: "Send a WhatsApp message to a specified phone number or group JID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"recipient": map[string]interface{}{
					"type":        "string",
					"description": "Phone number (with country code) or JID of the recipient",
				},
				"message": map[string]interface{}{
					"type":        "string",
					"description": "The message text to send",
				},
			},
			Required: []string{"recipient", "message"},
		},
	}, s.handleSendMessage)
}

// Tool handlers

func (s *Server) handleSearchContacts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Query string `json:"query"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	contacts, err := s.messenger.SearchContacts(ctx, args.Query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	result, _ := json.Marshal(contacts)
	return mcp.NewToolResultText(string(result)), nil
}

func (s *Server) handleListMessages(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		After     string `json:"after"`
		Before    string `json:"before"`
		SenderJID string `json:"sender_jid"`
		ChatJID   string `json:"chat_jid"`
		Query     string `json:"query"`
		Limit     int    `json:"limit"`
		Page      int    `json:"page"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	if args.Limit == 0 {
		args.Limit = 20
	}

	filter := messenger.MessageFilter{
		SenderJID: args.SenderJID,
		ChatJID:   args.ChatJID,
		Query:     args.Query,
		Limit:     args.Limit,
		Page:      args.Page,
	}

	if args.After != "" {
		t, err := time.Parse(time.RFC3339, args.After)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid after date: %v", err)), nil
		}
		filter.After = &t
	}

	if args.Before != "" {
		t, err := time.Parse(time.RFC3339, args.Before)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid before date: %v", err)), nil
		}
		filter.Before = &t
	}

	messages, err := s.messenger.ListMessages(ctx, filter)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("list messages failed: %v", err)), nil
	}

	result, _ := json.Marshal(messages)
	return mcp.NewToolResultText(string(result)), nil
}

func (s *Server) handleListChats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	if args.Limit == 0 {
		args.Limit = 20
	}

	chats, err := s.messenger.ListChats(ctx, args.Limit, args.Page)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("list chats failed: %v", err)), nil
	}

	result, _ := json.Marshal(chats)
	return mcp.NewToolResultText(string(result)), nil
}

func (s *Server) handleGetChat(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ChatJID string `json:"chat_jid"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	chat, err := s.messenger.GetChat(ctx, args.ChatJID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("get chat failed: %v", err)), nil
	}

	result, _ := json.Marshal(chat)
	return mcp.NewToolResultText(string(result)), nil
}

func (s *Server) handleGetDirectChatByContact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		PhoneNumber string `json:"phone_number"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	chat, err := s.messenger.GetDirectChatByContact(ctx, args.PhoneNumber)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("get direct chat failed: %v", err)), nil
	}

	result, _ := json.Marshal(chat)
	return mcp.NewToolResultText(string(result)), nil
}

func (s *Server) handleGetContactChats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ContactJID string `json:"contact_jid"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	chats, err := s.messenger.GetContactChats(ctx, args.ContactJID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("get contact chats failed: %v", err)), nil
	}

	result, _ := json.Marshal(chats)
	return mcp.NewToolResultText(string(result)), nil
}

func (s *Server) handleSendMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Recipient string `json:"recipient"`
		Message   string `json:"message"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	err := s.messenger.SendMessage(ctx, args.Recipient, args.Message)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("send message failed: %v", err)), nil
	}

	return mcp.NewToolResultText("Message sent successfully"), nil
}

// Serve starts the MCP server
func (s *Server) Serve() error {
	return server.ServeStdio(s.mcpServer)
}
