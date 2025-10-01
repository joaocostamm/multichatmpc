package whatsapp

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	"github.com/joao-costa/multichatmcp/internal/messenger"
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
				fmt.Printf("QR code: %s\n", evt.Code)
				log.Info().Msg("Scan the QR code above to log in")
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

// SearchContacts searches for contacts by name or phone number
func (w *WhatsAppMessenger) SearchContacts(ctx context.Context, query string) ([]messenger.Contact, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	contacts, err := w.client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	var results []messenger.Contact
	query = strings.ToLower(query)

	for jid, contact := range contacts {
		name := strings.ToLower(contact.FullName)
		phone := jid.User

		if strings.Contains(name, query) || strings.Contains(phone, query) {
			results = append(results, messenger.Contact{
				JID:         jid.String(),
				PhoneNumber: phone,
				Name:        contact.FullName,
			})
		}
	}

	return results, nil
}

// ListMessages retrieves messages with optional filters
func (w *WhatsAppMessenger) ListMessages(ctx context.Context, filter messenger.MessageFilter) ([]messenger.Message, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	// Note: whatsmeow doesn't provide direct message history access
	// You would need to implement message caching/storage separately
	// This is a placeholder implementation
	log.Warn().Msg("ListMessages: message history not fully implemented - requires custom message storage")

	return []messenger.Message{}, nil
}

// ListChats lists available chats
func (w *WhatsAppMessenger) ListChats(ctx context.Context, limit, page int) ([]messenger.Chat, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	// Get all contacts as a proxy for chats
	contacts, err := w.client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	var chats []messenger.Chat
	for jid, contact := range contacts {
		chat := messenger.Chat{
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
		return []messenger.Chat{}, nil
	}

	end := start + limit
	if end > len(chats) {
		end = len(chats)
	}

	return chats[start:end], nil
}

// GetChat gets information about a specific chat
func (w *WhatsAppMessenger) GetChat(ctx context.Context, chatJID string) (*messenger.Chat, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	jid, err := types.ParseJID(chatJID)
	if err != nil {
		return nil, fmt.Errorf("invalid JID: %w", err)
	}

	chat := &messenger.Chat{
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

// GetDirectChatByContact finds a direct chat with a specific contact
func (w *WhatsAppMessenger) GetDirectChatByContact(ctx context.Context, phoneNumber string) (*messenger.Chat, error) {
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
	return w.GetChat(ctx, jid.String())
}

// GetContactChats lists all chats involving a specific contact
func (w *WhatsAppMessenger) GetContactChats(ctx context.Context, contactJID string) ([]messenger.Chat, error) {
	if !w.IsConnected() {
		return nil, fmt.Errorf("not connected to WhatsApp")
	}

	// For direct messages, just return the direct chat
	chat, err := w.GetChat(ctx, contactJID)
	if err != nil {
		return nil, err
	}

	return []messenger.Chat{*chat}, nil
}

// SendMessage sends a message to a chat
func (w *WhatsAppMessenger) SendMessage(ctx context.Context, recipient, message string) error {
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
