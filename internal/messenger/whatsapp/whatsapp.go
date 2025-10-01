package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	_ "github.com/mattn/go-sqlite3"
	qrterminal "github.com/mdp/qrterminal/v3"
	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// WhatsAppMessenger implements the Messenger interface for WhatsApp
type WhatsAppMessenger struct {
	client    *whatsmeow.Client
	container *sqlstore.Container
	deviceDB  string
}

// NewWhatsAppMessenger creates a new WhatsApp messenger instance
func NewWhatsAppMessenger(deviceDB string) (*WhatsAppMessenger, error) {
	// Create a simple logger adapter for whatsmeow
	waLogger := waLog.Stdout("WhatsApp", "INFO", true)

	container, err := sqlstore.New(context.Background(), "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", deviceDB), waLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite container: %w", err)
	}

	return &WhatsAppMessenger{
		container: container,
		deviceDB:  deviceDB,
	}, nil
}

// Connect establishes connection to WhatsApp
func (w *WhatsAppMessenger) Connect(ctx context.Context) error {
	deviceStore, err := w.container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if deviceStore == nil {
		log.Info().Msg("No device found, creating new device")
		deviceStore = w.container.NewDevice()
	}

	w.client = whatsmeow.NewClient(deviceStore, nil)

	if w.client.Store.ID == nil {
		// No ID stored, new login required
		qrChan, err := w.client.GetQRChannel(ctx)
		if err != nil {
			return fmt.Errorf("failed to get QR channel: %w", err)
		}

		err = w.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				log.Info().Msg("QR code received, displaying for WhatsApp scan...")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				log.Info().Msg("Scan the QR code above with WhatsApp to log in")
			} else {
				log.Info().Str("event", evt.Event).Msg("Login event")
			}
		}
	} else {
		err = w.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
	}

	log.Info().Msg("WhatsApp connected successfully")
	return nil
}

// Disconnect closes the WhatsApp connection
func (w *WhatsAppMessenger) Disconnect() error {
	if w.client != nil {
		w.client.Disconnect()
	}
	if w.container != nil {
		return w.container.Close()
	}
	return nil
}

// GetMessengerName returns the name of the messenger platform
func (w *WhatsAppMessenger) GetMessengerName() string {
	return "whatsapp"
}

// SearchContacts searches for contacts by name or phone number
func (w *WhatsAppMessenger) searchContacts(ctx context.Context, query string) ([]Contact, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	contacts, err := w.client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	var results []Contact
	query = strings.ToLower(query)

	for jid, contact := range contacts {
		name := strings.ToLower(contact.FullName)
		phone := jid.User

		if strings.Contains(name, query) || strings.Contains(phone, query) {
			results = append(results, Contact{
				JID:         jid.String(),
				PhoneNumber: phone,
				Name:        contact.FullName,
			})
		}
	}

	return results, nil
}

// listMessages retrieves messages with optional filters
func (w *WhatsAppMessenger) listMessages(ctx context.Context, filter MessageFilter) ([]Message, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	// Note: whatsmeow doesn't provide direct message history access
	// You would need to implement message caching/storage separately
	// This is a placeholder implementation
	log.Warn().Msg("ListMessages: message history not fully implemented - requires custom message storage")

	return []Message{}, nil
}

// listChats lists available chats
func (w *WhatsAppMessenger) listChats(ctx context.Context, limit, page int) ([]Chat, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	// Get all contacts as a proxy for chats
	contacts, err := w.client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	var chats []Chat
	for jid, contact := range contacts {
		chat := Chat{
			JID:     jid.String(),
			IsGroup: jid.Server == types.GroupServer,
			Name:    contact.FullName,
		}

		if chat.Name == "" {
			chat.Name = jid.User
		}

		chats = append(chats, chat)
	}

	// Apply pagination
	start := page * limit
	if start >= len(chats) {
		return []Chat{}, nil
	}

	end := start + limit
	if end > len(chats) {
		end = len(chats)
	}

	return chats[start:end], nil
}

// getChat gets information about a specific chat
func (w *WhatsAppMessenger) getChat(ctx context.Context, chatJID string) (*Chat, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	jid, err := types.ParseJID(chatJID)
	if err != nil {
		return nil, fmt.Errorf("invalid JID: %w", err)
	}

	chat := &Chat{
		JID:     jid.String(),
		IsGroup: jid.Server == types.GroupServer,
	}

	if contact, err := w.client.Store.Contacts.GetContact(ctx, jid); err == nil {
		chat.Name = contact.FullName
	}
	if chat.Name == "" {
		chat.Name = jid.User
	}

	return chat, nil
}

// getDirectChatByContact finds a direct chat with a specific contact
func (w *WhatsAppMessenger) getDirectChatByContact(ctx context.Context, phoneNumber string) (*Chat, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	// Remove any non-numeric characters
	phone := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phoneNumber)

	jid := types.NewJID(phone, types.DefaultUserServer)
	return w.getChat(ctx, jid.String())
}

// getContactChats lists all chats involving a specific contact
func (w *WhatsAppMessenger) getContactChats(ctx context.Context, contactJID string) ([]Chat, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	// For direct messages, just return the direct chat
	chat, err := w.getChat(ctx, contactJID)
	if err != nil {
		return nil, err
	}

	return []Chat{*chat}, nil
}

// sendMessage sends a message to a chat
func (w *WhatsAppMessenger) sendMessage(ctx context.Context, recipient, message string) error {
	if !w.IsConnected() {
		return fmt.Errorf("not connected to WhatsApp")
	}

	// Parse recipient as JID or phone number
	var jid types.JID
	var err error

	if strings.Contains(recipient, "@") {
		jid, err = types.ParseJID(recipient)
		if err != nil {
			return fmt.Errorf("invalid JID: %w", err)
		}
	} else {
		// Assume it's a phone number
		phone := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, recipient)
		jid = types.NewJID(phone, types.DefaultUserServer)
	}

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err = w.client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Info().Str("recipient", jid.String()).Msg("Message sent")
	return nil
}

// IsConnected returns the connection status
func (w *WhatsAppMessenger) IsConnected() bool {
	return w.client != nil && w.client.IsConnected()
}

// RegisterMCPTools registers WhatsApp-specific MCP tools
func (w *WhatsAppMessenger) RegisterMCPTools(mcpServer *server.MCPServer) {
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
	}, w.handleSearchContacts)

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
	}, w.handleListMessages)

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
	}, w.handleListChats)

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
	}, w.handleGetChat)

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
	}, w.handleGetDirectChatByContact)

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
	}, w.handleGetContactChats)

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
	}, w.handleSendMessage)
}

// Tool handlers

func (w *WhatsAppMessenger) handleSearchContacts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Query string `json:"query"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	contacts, err := w.searchContacts(ctx, args.Query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	result, _ := json.Marshal(contacts)
	return mcp.NewToolResultText(string(result)), nil
}

func (w *WhatsAppMessenger) handleListMessages(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	filter := MessageFilter{
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

	messages, err := w.listMessages(ctx, filter)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("list messages failed: %v", err)), nil
	}

	result, _ := json.Marshal(messages)
	return mcp.NewToolResultText(string(result)), nil
}

func (w *WhatsAppMessenger) handleListChats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	chats, err := w.listChats(ctx, args.Limit, args.Page)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("list chats failed: %v", err)), nil
	}

	result, _ := json.Marshal(chats)
	return mcp.NewToolResultText(string(result)), nil
}

func (w *WhatsAppMessenger) handleGetChat(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ChatJID string `json:"chat_jid"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	chat, err := w.getChat(ctx, args.ChatJID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("get chat failed: %v", err)), nil
	}

	result, _ := json.Marshal(chat)
	return mcp.NewToolResultText(string(result)), nil
}

func (w *WhatsAppMessenger) handleGetDirectChatByContact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		PhoneNumber string `json:"phone_number"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	chat, err := w.getDirectChatByContact(ctx, args.PhoneNumber)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("get direct chat failed: %v", err)), nil
	}

	result, _ := json.Marshal(chat)
	return mcp.NewToolResultText(string(result)), nil
}

func (w *WhatsAppMessenger) handleGetContactChats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ContactJID string `json:"contact_jid"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	chats, err := w.getContactChats(ctx, args.ContactJID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("get contact chats failed: %v", err)), nil
	}

	result, _ := json.Marshal(chats)
	return mcp.NewToolResultText(string(result)), nil
}

func (w *WhatsAppMessenger) handleSendMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Recipient string `json:"recipient"`
		Message   string `json:"message"`
	}
	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	err := w.sendMessage(ctx, args.Recipient, args.Message)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("send message failed: %v", err)), nil
	}

	return mcp.NewToolResultText("Message sent successfully"), nil
}
