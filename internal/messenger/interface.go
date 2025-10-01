package messenger

import (
	"context"
	"time"
)

// Contact represents a messenger contact
type Contact struct {
	JID         string `json:"jid"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	ChatJID   string    `json:"chat_jid"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	IsFromMe  bool      `json:"is_from_me"`
	MediaType string    `json:"media_type,omitempty"`
}

// Chat represents a conversation
type Chat struct {
	JID         string   `json:"jid"`
	Name        string   `json:"name"`
	IsGroup     bool     `json:"is_group"`
	LastMessage *Message `json:"last_message,omitempty"`
}

// MessageFilter contains criteria for filtering messages
type MessageFilter struct {
	After     *time.Time
	Before    *time.Time
	SenderJID string
	ChatJID   string
	Query     string
	Limit     int
	Page      int
}

// Messenger defines the interface that all messenger implementations must follow
type Messenger interface {
	// Connect establishes connection to the messenger service
	Connect(ctx context.Context) error

	// Disconnect closes the connection
	Disconnect() error

	// SearchContacts searches for contacts by name or phone number
	SearchContacts(ctx context.Context, query string) ([]Contact, error)

	// ListMessages retrieves messages with optional filters
	ListMessages(ctx context.Context, filter MessageFilter) ([]Message, error)

	// ListChats lists available chats
	ListChats(ctx context.Context, limit, page int) ([]Chat, error)

	// GetChat gets information about a specific chat
	GetChat(ctx context.Context, chatJID string) (*Chat, error)

	// GetDirectChatByContact finds a direct chat with a specific contact
	GetDirectChatByContact(ctx context.Context, phoneNumber string) (*Chat, error)

	// GetContactChats lists all chats involving a specific contact
	GetContactChats(ctx context.Context, contactJID string) ([]Chat, error)

	// SendMessage sends a message to a chat
	SendMessage(ctx context.Context, recipient, message string) error

	// IsConnected returns the connection status
	IsConnected() bool
}
