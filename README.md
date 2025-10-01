# MultiChat MCP Server

A multi-messenger MCP (Model Context Protocol) server implementation in Go. Currently supports WhatsApp with a modular architecture designed to support additional messaging platforms in the future.

## Features

- ✅ **WhatsApp Support** using [whatsmeow](https://github.com/tulir/whatsmeow)
- ✅ **MCP Protocol** implementation using [mcp-go](https://github.com/modelcontextprotocol/go-mcp)
- ✅ **CLI Interface** with [Cobra](https://github.com/spf13/cobra)
- ✅ **Structured Logging** with [zerolog](https://github.com/rs/zerolog)
- ✅ **Modular Architecture** for easy addition of new messaging platforms

## Installation

### Prerequisites

- Go 1.25.1 or higher
- SQLite3

### Build from Source

```bash
# Clone the repository
git clone https://github.com/joao-costa/multichatmcp.git
cd multichatmcp

# Download dependencies
go mod download

# Build the application
make build

# Or build directly with go
go build -o multichat

# Or install to GOPATH/bin
make install
```

## Usage

### First-Time Setup (WhatsApp)

When running for the first time, the application will generate a QR code for WhatsApp authentication:

```bash
./multichat --messenger whatsapp --device mydevice.db --log-level debug
```

1. Scan the QR code with your WhatsApp mobile app (Settings → Linked Devices → Link a Device)
2. The session will be saved to `mydevice.db` for future use
3. Subsequent runs will automatically reconnect using the saved session

### Command-Line Arguments

```bash
./multichat [flags]

Flags:
  --messenger string    Messenger type (default "whatsapp")
  --device string       Device database file path (default "device.db")
  --log-level string    Log level: debug, info, warn, error (default "info")
  -h, --help           Help for multichat
```

### MCP Server Configuration

To use with an MCP client (like Claude Desktop, Cursor, or other AI assistants), add the following to your MCP configuration:

#### Example: `mcp_config.json`

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "/path/to/multichat",
      "args": [
        "--messenger",
        "whatsapp",
        "--device",
        "mydevice.db"
      ]
    }
  }
}
```

#### For Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS):

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "/Users/yourusername/path/to/multichat",
      "args": [
        "--messenger",
        "whatsapp",
        "--device",
        "/Users/yourusername/path/to/mydevice.db",
        "--log-level",
        "info"
      ]
    }
  }
}
```

## MCP Operations

The server exposes the following MCP tools:

### 1. `search_contacts`

Search for contacts by name or phone number.

**Parameters:**
- `query` (string, required): Search term to match against contact names or phone numbers

**Example:**
```json
{
  "query": "John"
}
```

### 2. `list_messages`

Retrieve messages with optional filters and context.

**Parameters:**
- `after` (string, optional): ISO-8601 formatted date to only return messages after this date
- `before` (string, optional): ISO-8601 formatted date to only return messages before this date
- `sender_jid` (string, optional): Filter messages by sender JID
- `chat_jid` (string, optional): Filter messages by chat JID
- `query` (string, optional): Search term to filter messages by content
- `limit` (integer, optional): Maximum number of messages to return (default: 20)
- `page` (integer, optional): Page number for pagination (default: 0)

**Example:**
```json
{
  "chat_jid": "1234567890@s.whatsapp.net",
  "limit": 50
}
```

**Note:** Full message history retrieval requires implementing a custom message storage layer, as whatsmeow doesn't provide direct access to message history.

### 3. `list_chats`

List available chats with metadata.

**Parameters:**
- `limit` (integer, optional): Maximum number of chats to return (default: 20)
- `page` (integer, optional): Page number for pagination (default: 0)

**Example:**
```json
{
  "limit": 10,
  "page": 0
}
```

### 4. `get_chat`

Get information about a specific chat.

**Parameters:**
- `chat_jid` (string, required): The JID of the chat to retrieve

**Example:**
```json
{
  "chat_jid": "1234567890@s.whatsapp.net"
}
```

### 5. `get_direct_chat_by_contact`

Find a direct chat with a specific contact by phone number.

**Parameters:**
- `phone_number` (string, required): Phone number of the contact (with country code, no + or spaces)

**Example:**
```json
{
  "phone_number": "1234567890"
}
```

### 6. `get_contact_chats`

List all chats involving a specific contact.

**Parameters:**
- `contact_jid` (string, required): The JID of the contact

**Example:**
```json
{
  "contact_jid": "1234567890@s.whatsapp.net"
}
```

### 7. `send_message`

Send a WhatsApp message to a specified phone number or group JID.

**Parameters:**
- `recipient` (string, required): Phone number (with country code) or JID of the recipient
- `message` (string, required): The message text to send

**Example (phone number):**
```json
{
  "recipient": "1234567890",
  "message": "Hello from MCP!"
}
```

**Example (JID):**
```json
{
  "recipient": "1234567890@s.whatsapp.net",
  "message": "Hello from MCP!"
}
```

**Example (group):**
```json
{
  "recipient": "1234567890@g.us",
  "message": "Hello group!"
}
```

## Project Structure

```
multichatmcp/
├── main.go                 # Main application entry point
├── internal/
│   ├── messenger/          # Messenger interface and implementations
│   │   ├── interface.go    # Messenger interface definition
│   │   └── whatsapp/       # WhatsApp implementation
│   │       └── whatsapp.go
│   └── mcp/                # MCP server implementation
│       └── server.go
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
└── README.md
```

## Architecture

The application follows a modular architecture:

1. **Messenger Interface**: Defines a common interface that all messaging platforms must implement
2. **Platform Implementations**: Each messaging platform (WhatsApp, future: Telegram, Signal, etc.) implements the Messenger interface
3. **MCP Server**: Wraps the messenger implementation and exposes MCP tools
4. **CLI**: Cobra-based CLI for configuration and startup

This design makes it easy to add support for new messaging platforms by implementing the `Messenger` interface.

## Adding New Messengers

To add support for a new messaging platform:

1. Create a new package under `internal/messenger/<platform>/`
2. Implement the `messenger.Messenger` interface
3. Add the new messenger type to the CLI switch in `main.go`
4. Update documentation

Example interface methods to implement:
- `Connect(ctx context.Context) error`
- `Disconnect() error`
- `SearchContacts(ctx context.Context, query string) ([]Contact, error)`
- `ListMessages(ctx context.Context, filter MessageFilter) ([]Message, error)`
- `ListChats(ctx context.Context, limit, page int) ([]Chat, error)`
- `GetChat(ctx context.Context, chatJID string) (*Chat, error)`
- `GetDirectChatByContact(ctx context.Context, phoneNumber string) (*Chat, error)`
- `GetContactChats(ctx context.Context, contactJID string) ([]Chat, error)`
- `SendMessage(ctx context.Context, recipient, message string) error`
- `IsConnected() bool`

## Limitations & Known Issues

1. **Message History**: WhatsApp's `whatsmeow` library doesn't provide direct access to message history. To implement full message retrieval, you would need to:
   - Set up event handlers to capture incoming messages
   - Store messages in a local database
   - Query from your local database in the `ListMessages` implementation

2. **Media Messages**: Currently only text messages are supported. Media support can be added by extending the implementation.

3. **Group Management**: Advanced group operations (create, manage members, etc.) are not yet implemented.

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
# Build with optimizations
make build-prod

# Or directly with go
go build -ldflags="-s -w" -o multichat

# Cross-compile for all platforms
make build-all
```

## Troubleshooting

### QR Code Not Appearing

- Ensure you're running with `--log-level debug` to see detailed logs
- Check that no other WhatsApp Web/Desktop sessions are active
- Delete the device database file and try again

### Connection Issues

- Verify your internet connection
- Check firewall settings
- Ensure WhatsApp is working on your mobile device

### Database Locked Errors

- Make sure only one instance of the application is running
- Check file permissions on the device database

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [whatsmeow](https://github.com/tulir/whatsmeow) - WhatsApp Web multidevice API library
- [mcp-go](https://github.com/modelcontextprotocol/go-mcp) - Model Context Protocol implementation for Go
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [zerolog](https://github.com/rs/zerolog) - Zero-allocation JSON logger

