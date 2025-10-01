package teams

// TeamsConfig holds configuration for Teams messenger
type TeamsConfig struct {
	// DefaultWebhookURL is the default webhook URL to use if not specified per message
	DefaultWebhookURL string `json:"default_webhook_url,omitempty"`
}

// MessageCard represents the response from sending a message
type MessageCard struct {
	WebhookURL string `json:"webhook_url"`
	Title      string `json:"title,omitempty"`
	Text       string `json:"text"`
	Color      string `json:"color,omitempty"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}
