package whatsapp

import "time"

// Contact represents a WhatsApp contact
type Contact struct {
	JID         string `json:"jid"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

// Message represents a WhatsApp chat message
type Message struct {
	ID        string    `json:"id"`
	ChatJID   string    `json:"chat_jid"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	IsFromMe  bool      `json:"is_from_me"`
	MediaType string    `json:"media_type,omitempty"`
}

// Chat represents a WhatsApp conversation
type Chat struct {
	JID         string   `json:"jid"`
	Name        string   `json:"name"`
	IsGroup     bool     `json:"is_group"`
	LastMessage *Message `json:"last_message,omitempty"`
}

// MessageFilter contains criteria for filtering WhatsApp messages
type MessageFilter struct {
	After     *time.Time
	Before    *time.Time
	SenderJID string
	ChatJID   string
	Query     string
	Limit     int
	Page      int
}
