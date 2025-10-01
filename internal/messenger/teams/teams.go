package teams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// TeamsMessenger implements the Messenger interface for Microsoft Teams
type TeamsMessenger struct {
	config    TeamsConfig
	client    *goteamsnotify.TeamsClient
	connected bool
}

// NewTeamsMessenger creates a new Teams messenger instance
func NewTeamsMessenger(webhookURL string) (*TeamsMessenger, error) {
	return &TeamsMessenger{
		config: TeamsConfig{
			DefaultWebhookURL: webhookURL,
		},
		client:    goteamsnotify.NewTeamsClient(),
		connected: false,
	}, nil
}

// Connect validates the webhook URL and prepares the messenger
func (t *TeamsMessenger) Connect(ctx context.Context) error {
	// Validate default webhook URL if provided
	if t.config.DefaultWebhookURL != "" {
		if err := t.validateWebhookURL(t.config.DefaultWebhookURL); err != nil {
			return fmt.Errorf("invalid default webhook URL: %w", err)
		}
		log.Info().Msg("Default webhook URL validated successfully")
	} else {
		log.Warn().Msg("No default webhook URL provided - webhook URL must be specified for each message")
	}

	t.connected = true
	log.Info().Msg("Teams messenger connected successfully")
	return nil
}

// Disconnect closes any resources (none needed for Teams)
func (t *TeamsMessenger) Disconnect() error {
	t.connected = false
	log.Info().Msg("Teams messenger disconnected")
	return nil
}

// IsConnected returns the connection status
func (t *TeamsMessenger) IsConnected() bool {
	return t.connected
}

// GetMessengerName returns the name of the messenger platform
func (t *TeamsMessenger) GetMessengerName() string {
	return "teams"
}

// validateWebhookURL validates that the URL is a valid Teams webhook URL
func (t *TeamsMessenger) validateWebhookURL(webhookURL string) error {
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if it's a valid Teams webhook URL (Power Automate or O365 connector)
	// Power Automate URLs typically contain "prod.apiflow.microsoft.com" or similar
	// O365 connector URLs contain "outlook.office.com" or "webhook.office.com"
	validHosts := []string{
		"prod.apiflow.microsoft.com",
		"prod-",
		"outlook.office.com",
		"outlook.office365.com",
		"webhook.office.com",
	}

	isValid := false
	for _, validHost := range validHosts {
		if parsedURL.Host == validHost || parsedURL.Hostname() == validHost {
			isValid = true
			break
		}
		// Also check if host contains the valid pattern (for regional endpoints)
		if len(parsedURL.Host) > len(validHost) && parsedURL.Host[:len(validHost)] == validHost {
			isValid = true
			break
		}
	}

	if !isValid {
		log.Warn().Str("host", parsedURL.Host).Msg("Webhook URL host doesn't match known Teams patterns - will attempt anyway")
	}

	return nil
}

// sendMessage sends a message to a Teams channel or chat via webhook
func (t *TeamsMessenger) sendMessage(ctx context.Context, webhookURL, message, title, color string) (*MessageCard, error) {
	if !t.IsConnected() {
		return nil, fmt.Errorf("not connected to Teams")
	}

	// Use default webhook URL if not specified
	if webhookURL == "" {
		webhookURL = t.config.DefaultWebhookURL
	}

	if webhookURL == "" {
		return nil, fmt.Errorf("webhook URL is required - either provide it or set a default webhook URL")
	}

	// Validate webhook URL
	if err := t.validateWebhookURL(webhookURL); err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Create message card
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = title
	msgCard.Text = message

	// Set color if provided
	if color != "" {
		msgCard.ThemeColor = color
	}

	// Send the message
	if err := t.client.Send(webhookURL, msgCard); err != nil {
		log.Error().Err(err).Str("webhook", webhookURL).Msg("Failed to send Teams message")
		return &MessageCard{
			WebhookURL: webhookURL,
			Title:      title,
			Text:       message,
			Color:      color,
			Success:    false,
			Error:      err.Error(),
		}, err
	}

	log.Info().Str("webhook", webhookURL).Msg("Teams message sent successfully")
	return &MessageCard{
		WebhookURL: webhookURL,
		Title:      title,
		Text:       message,
		Color:      color,
		Success:    true,
	}, nil
}

// sendRichMessage sends a rich message card with facts
func (t *TeamsMessenger) sendRichMessage(ctx context.Context, webhookURL, title, text, color string, facts map[string]string) (*MessageCard, error) {
	if !t.IsConnected() {
		return nil, fmt.Errorf("not connected to Teams")
	}

	// Use default webhook URL if not specified
	if webhookURL == "" {
		webhookURL = t.config.DefaultWebhookURL
	}

	if webhookURL == "" {
		return nil, fmt.Errorf("webhook URL is required")
	}

	// Validate webhook URL
	if err := t.validateWebhookURL(webhookURL); err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Create message card
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = title
	msgCard.Text = text

	// Set theme color
	if color != "" {
		msgCard.ThemeColor = color
	}

	// Add facts as a section if provided
	if len(facts) > 0 {
		section := messagecard.NewSection()
		section.Title = "Details"
		for key, value := range facts {
			section.AddFactFromKeyValue(key, value)
		}
		if err := msgCard.AddSection(section); err != nil {
			return nil, fmt.Errorf("failed to add facts section: %w", err)
		}
	}

	// Send the message
	if err := t.client.Send(webhookURL, msgCard); err != nil {
		log.Error().Err(err).Str("webhook", webhookURL).Msg("Failed to send Teams rich message")
		return &MessageCard{
			WebhookURL: webhookURL,
			Title:      title,
			Text:       text,
			Color:      color,
			Success:    false,
			Error:      err.Error(),
		}, err
	}

	log.Info().Str("webhook", webhookURL).Msg("Teams rich message sent successfully")
	return &MessageCard{
		WebhookURL: webhookURL,
		Title:      title,
		Text:       text,
		Color:      color,
		Success:    true,
	}, nil
}

// RegisterMCPTools registers Teams-specific MCP tools
func (t *TeamsMessenger) RegisterMCPTools(mcpServer *server.MCPServer) {
	// send_message - Simple text message
	mcpServer.AddTool(mcp.Tool{
		Name:        "send_message",
		Description: "Send a simple message to a Microsoft Teams channel or chat via webhook URL",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"webhook_url": map[string]interface{}{
					"type":        "string",
					"description": "Teams webhook URL (Power Automate workflow URL or O365 connector URL). If not provided, uses the default webhook URL set at initialization.",
				},
				"message": map[string]interface{}{
					"type":        "string",
					"description": "The message text to send",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Optional title for the message card",
				},
				"color": map[string]interface{}{
					"type":        "string",
					"description": "Optional theme color in hex format (e.g., '0078D4' for blue, 'FF0000' for red)",
				},
			},
			Required: []string{"message"},
		},
	}, t.handleSendMessage)

	// send_rich_message - Rich message with facts
	mcpServer.AddTool(mcp.Tool{
		Name:        "send_rich_message",
		Description: "Send a rich adaptive card message with title, text, color, and structured facts to Teams",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"webhook_url": map[string]interface{}{
					"type":        "string",
					"description": "Teams webhook URL. If not provided, uses the default webhook URL.",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Title of the message card",
				},
				"text": map[string]interface{}{
					"type":        "string",
					"description": "Main text content of the message",
				},
				"color": map[string]interface{}{
					"type":        "string",
					"description": "Theme color in hex format (e.g., '0078D4', 'FF0000', '00FF00')",
				},
				"facts": map[string]interface{}{
					"type":        "object",
					"description": "Key-value pairs to display as facts (e.g., {'Status': 'Active', 'Priority': 'High'})",
				},
			},
			Required: []string{"text"},
		},
	}, t.handleSendRichMessage)

	// validate_webhook - Validate a webhook URL
	mcpServer.AddTool(mcp.Tool{
		Name:        "validate_webhook",
		Description: "Validate a Teams webhook URL to ensure it's properly formatted",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"webhook_url": map[string]interface{}{
					"type":        "string",
					"description": "Teams webhook URL to validate",
				},
			},
			Required: []string{"webhook_url"},
		},
	}, t.handleValidateWebhook)
}

// Tool handlers

func (t *TeamsMessenger) handleSendMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		WebhookURL string `json:"webhook_url"`
		Message    string `json:"message"`
		Title      string `json:"title"`
		Color      string `json:"color"`
	}

	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	result, err := t.sendMessage(ctx, args.WebhookURL, args.Message, args.Title, args.Color)
	if err != nil {
		// Return the error details in the result
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultError(fmt.Sprintf("failed to send message: %v\nDetails: %s", err, string(resultJSON))), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (t *TeamsMessenger) handleSendRichMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		WebhookURL string            `json:"webhook_url"`
		Title      string            `json:"title"`
		Text       string            `json:"text"`
		Color      string            `json:"color"`
		Facts      map[string]string `json:"facts"`
	}

	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	result, err := t.sendRichMessage(ctx, args.WebhookURL, args.Title, args.Text, args.Color, args.Facts)
	if err != nil {
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultError(fmt.Sprintf("failed to send rich message: %v\nDetails: %s", err, string(resultJSON))), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (t *TeamsMessenger) handleValidateWebhook(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		WebhookURL string `json:"webhook_url"`
	}

	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	if err := t.validateWebhookURL(args.WebhookURL); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf(`{"valid": false, "webhook_url": "%s", "error": "%v"}`, args.WebhookURL, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"valid": true, "webhook_url": "%s"}`, args.WebhookURL)), nil
}
