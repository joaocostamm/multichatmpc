<div align="center">

# 🚀 MultiChat MCP Server

**Bridge your messaging platforms with AI through the Model Context Protocol**

[![Go Version](https://img.shields.io/badge/Go-1.25.1+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-blue?style=for-the-badge)](https://modelcontextprotocol.io)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=for-the-badge)](CONTRIBUTING.md)

[Features](#-features) •
[Quick Start](#-quick-start) •
[Documentation](#-documentation) •
[Architecture](#-architecture) •
[Contributing](#-contributing)

</div>

---

## 🎯 What is MultiChat MCP?

MultiChat MCP Server is a powerful Go-based implementation of the [Model Context Protocol (MCP)](https://modelcontextprotocol.io) that enables AI assistants like Claude, GPT, and others to seamlessly interact with your messaging platforms. Start with WhatsApp and expand to Teams, Telegram, Signal, and more with our **messenger-specific operations architecture**.

### Why MultiChat MCP?

- 🤖 **AI-Native**: Built specifically for AI assistants to read and send messages
- 🔌 **Plug & Play**: Easy integration with Claude Desktop, Cursor, and any MCP-compatible client
- 🧩 **Truly Modular**: Each messenger defines its own operations - no forced common interface
- 🎯 **Platform-Specific**: Each messenger exposes only the operations that make sense for that platform
- 🔒 **Privacy First**: Your data stays on your machine
- ⚡ **Lightning Fast**: Written in Go for optimal performance
- 🔧 **Dynamic**: MCP tools are registered at runtime based on the selected messenger

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 📱 Platform Support
- ✅ **WhatsApp** - 7 operations (via [whatsmeow](https://github.com/tulir/whatsmeow))
- ✅ **Teams** - 3 operations (via [go-teams-notify](https://github.com/atc0005/go-teams-notify))
- ✅ **Twitter/X** - 3 operations (via [gotwi](https://github.com/michimani/gotwi))
- 🔜 **Telegram** - Platform-specific tools (polls, forwards, etc.)
- 🔜 **Signal** - Secure messaging operations
- 🔜 **Discord** - Server/channel management

*Each platform has its own unique set of MCP operations*

</td>
<td width="50%">

### 🛠️ Architecture Capabilities
- 🔧 **Messenger-specific operations** - Each platform defines its own tools
- 🔌 **Dynamic MCP registration** - Tools registered at runtime
- 🧩 **Minimal interface** - No forced abstractions
- 🎯 **Platform isolation** - Independent implementations
- 📦 **Type safety** - Platform-specific types

</td>
</tr>
</table>

### 🎨 Technology Stack

- **[MCP Protocol](https://github.com/modelcontextprotocol/go-mcp)** - Standard protocol for AI-app communication
- **[Cobra](https://github.com/spf13/cobra)** - Modern CLI framework
- **[Zerolog](https://github.com/rs/zerolog)** - Zero-allocation structured logging
- **[SQLite3](https://www.sqlite.org/)** - Lightweight session storage

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.25.1+** ([Download](https://golang.org/dl/))
- **SQLite3** (usually pre-installed on macOS/Linux)

### Installation

```bash
# Clone the repository
git clone https://github.com/joao-costa/multichatmcp.git
cd multichatmcp

# Build the binary
make build

# Or install directly to your GOPATH
make install
```

### First Run

#### WhatsApp Setup

```bash
./multichat --messenger whatsapp --device mydevice.db --log-level debug
```

**🔐 Authentication Steps:**
1. A QR code will appear in your terminal
2. Open WhatsApp on your phone → **Settings** → **Linked Devices** → **Link a Device**
3. Scan the QR code
4. Done! Your session is saved for future use

#### Teams Setup

```bash
./multichat --messenger teams --webhook "https://your-webhook-url" --log-level debug
```

**🔗 Webhook Setup Steps:**
1. Create a Power Automate Workflow webhook URL (see [Webhook Setup](#teams-webhook-setup))
2. Run the server with your webhook URL
3. The webhook URL can also be provided per-message for multi-channel support

#### Twitter/X Setup

```bash
./multichat --messenger twitter \
  --twitter-api-key "your-api-key" \
  --twitter-api-secret "your-api-secret" \
  --twitter-token "your-access-token" \
  --twitter-token-secret "your-access-token-secret" \
  --log-level debug
```

**🔑 Twitter API Credentials Setup:**
1. Go to [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard)
2. Create a new project and app (or use existing)
3. Navigate to your app's **Keys and Tokens** section
4. Generate or retrieve:
   - API Key (also called Consumer Key)
   - API Secret Key (also called Consumer Secret)
   - Access Token
   - Access Token Secret
5. Ensure your app has **Read and Write** permissions to post tweets
6. Use the credentials with the flags above

---

## 📖 Documentation

### Command-Line Usage

```bash
./multichat [flags]

Flags:
  --messenger string           Messaging platform to use: whatsapp, teams, twitter (default "whatsapp")
  --device string             Device database file path (for WhatsApp) (default "device.db")
  --webhook string            Webhook URL (for Teams) (optional, can be provided per-message)
  --twitter-api-key string    Twitter API Key (for Twitter/X)
  --twitter-api-secret string Twitter API Secret Key (for Twitter/X)
  --twitter-token string      Twitter Access Token (for Twitter/X)
  --twitter-token-secret string Twitter Access Token Secret (for Twitter/X)
  --log-level string          Logging level: debug, info, warn, error (default "info")
  -h, --help                  Show help information
```

### MCP Client Configuration

#### 🖥️ Claude Desktop

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`

**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "/absolute/path/to/multichat",
      "args": [
        "--messenger", "whatsapp",
        "--device", "/absolute/path/to/device.db",
        "--log-level", "info"
      ]
    },
    "teams": {
      "command": "/absolute/path/to/multichat",
      "args": [
        "--messenger", "teams",
        "--webhook", "https://your-teams-webhook-url",
        "--log-level", "info"
      ]
    },
    "twitter": {
      "command": "/absolute/path/to/multichat",
      "args": [
        "--messenger", "twitter",
        "--twitter-api-key", "your-api-key",
        "--twitter-api-secret", "your-api-secret",
        "--twitter-token", "your-access-token",
        "--twitter-token-secret", "your-access-token-secret",
        "--log-level", "info"
      ]
    }
  }
}
```

#### 🎯 Cursor IDE

**Location:** `~/.cursor/mcp.json` or workspace settings

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "/absolute/path/to/multichat",
      "args": ["--messenger", "whatsapp", "--device", "device.db"]
    }
  }
}
```

#### 🔧 Generic MCP Client

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "/path/to/multichat",
      "args": ["--messenger", "whatsapp", "--device", "mydevice.db"]
    }
  }
}
```

---

## 🔧 Available MCP Tools

### WhatsApp Tools

When running with `--messenger whatsapp`, the following MCP tools are available:

### 👤 `search_contacts`
Find contacts by name or phone number.

```json
{
  "query": "John Doe"
}
```

### 💬 `list_messages`
Retrieve messages with powerful filtering options.

```json
{
  "chat_jid": "1234567890@s.whatsapp.net",
  "after": "2024-01-01T00:00:00Z",
  "before": "2024-12-31T23:59:59Z",
  "query": "meeting",
  "limit": 50,
  "page": 0
}
```

**Parameters:**
- `after` *(string, optional)*: ISO-8601 date - messages after this time
- `before` *(string, optional)*: ISO-8601 date - messages before this time
- `sender_jid` *(string, optional)*: Filter by sender
- `chat_jid` *(string, optional)*: Filter by chat
- `query` *(string, optional)*: Full-text search
- `limit` *(integer, optional)*: Max results (default: 20)
- `page` *(integer, optional)*: Page number (default: 0)

### 📋 `list_chats`
Get all available chats with metadata.

```json
{
  "limit": 20,
  "page": 0
}
```

### 🔍 `get_chat`
Retrieve detailed information about a specific chat.

```json
{
  "chat_jid": "1234567890@s.whatsapp.net"
}
```

### 📞 `get_direct_chat_by_contact`
Find direct chat by phone number.

```json
{
  "phone_number": "1234567890"
}
```

**Note:** Use country code without `+` or spaces (e.g., `15551234567` for US)

### 👥 `get_contact_chats`
List all chats involving a specific contact.

```json
{
  "contact_jid": "1234567890@s.whatsapp.net"
}
```

### 📤 `send_message`
Send a message to any contact or group.

**Direct Message:**
```json
{
  "recipient": "1234567890",
  "message": "Hello from MultiChat MCP! 👋"
}
```

**Using JID:**
```json
{
  "recipient": "1234567890@s.whatsapp.net",
  "message": "Hey there!"
}
```

**Group Message:**
```json
{
  "recipient": "1234567890@g.us",
  "message": "Hello everyone! 🎉"
}
```

---

### Teams Tools

When running with `--messenger teams`, the following MCP tools are available:

#### 📤 `send_message`
Send a simple message to a Teams channel or chat via webhook.

```json
{
  "webhook_url": "https://your-webhook-url",
  "message": "Hello from MultiChat MCP! 🚀",
  "title": "Notification",
  "color": "0078D4"
}
```

**Parameters:**
- `webhook_url` *(string, optional)*: Teams webhook URL. If not provided, uses the default webhook URL set at initialization
- `message` *(string, required)*: The message text to send
- `title` *(string, optional)*: Optional title for the message card
- `color` *(string, optional)*: Theme color in hex format (e.g., '0078D4' for blue, 'FF0000' for red, '00FF00' for green)

#### 📊 `send_rich_message`
Send a rich message with title, text, color, and structured facts.

```json
{
  "webhook_url": "https://your-webhook-url",
  "title": "Deployment Status",
  "text": "The deployment has completed successfully",
  "color": "00FF00",
  "facts": {
    "Environment": "Production",
    "Version": "v1.2.3",
    "Status": "Success"
  }
}
```

**Parameters:**
- `webhook_url` *(string, optional)*: Teams webhook URL
- `title` *(string, optional)*: Title of the message card
- `text` *(string, required)*: Main text content
- `color` *(string, optional)*: Theme color in hex format
- `facts` *(object, optional)*: Key-value pairs to display as structured facts

#### ✅ `validate_webhook`
Validate a Teams webhook URL to ensure it's properly formatted.

```json
{
  "webhook_url": "https://your-webhook-url"
}
```

**Returns:** `{"valid": true/false, "webhook_url": "...", "error": "..."}`

---

### Twitter/X Tools

When running with `--messenger twitter`, the following MCP tools are available:

#### 🐦 `post_tweet`
Post a tweet to Twitter/X (max 280 characters).

```json
{
  "text": "Hello from MultiChat MCP! 🚀",
  "reply_to_tweet_id": "1234567890123456789"
}
```

**Parameters:**
- `text` *(string, required)*: The text content of the tweet (max 280 characters)
- `reply_to_tweet_id` *(string, optional)*: ID of the tweet to reply to (creates a reply thread)

**Returns:** `{"tweet_id": "...", "text": "...", "success": true}`

#### 📤 `send_message`
Send a message (tweet) to Twitter/X - alias for `post_tweet` to maintain consistency with other messengers.

```json
{
  "message": "Just posted via MCP! 🎉",
  "reply_to_tweet_id": "1234567890123456789"
}
```

**Parameters:**
- `message` *(string, required)*: The text content of the tweet (max 280 characters)
- `reply_to_tweet_id` *(string, optional)*: ID of the tweet to reply to

**Returns:** `{"tweet_id": "...", "text": "...", "success": true}`

#### 🗑️ `delete_tweet`
Delete a tweet by ID.

```json
{
  "tweet_id": "1234567890123456789"
}
```

**Parameters:**
- `tweet_id` *(string, required)*: The ID of the tweet to delete

**Returns:** `{"success": true, "tweet_id": "...", "message": "Tweet deleted successfully"}`

**Note:** You can only delete tweets posted by the authenticated account.

---

### Teams Webhook Setup

To use Teams integration, you need to create a Power Automate Workflow webhook URL. Microsoft has deprecated the old O365 connectors (retiring October 2024), so you must use Power Automate workflows.

#### Creating a Workflow Webhook (Power Automate)

**Method 1: Using Teams Client - Workflows Context Option**

1. Navigate to a channel or chat in Teams
2. Select the ellipsis (⋯) on the channel or chat
3. Select **Workflows**
4. Type `when a webhook request` in the search
5. Select **Post to a channel when a webhook request is received**
6. Verify that **Microsoft Teams** is successfully enabled
7. Select **Next**
8. Select an appropriate **Team** from the dropdown
9. Select an appropriate **Channel** from the dropdown
10. Select **Create flow**
11. Copy the new workflow URL (this is your webhook URL)
12. Select **Done**

**Method 2: Using Power Automate Web UI**

1. Go to [https://make.powerautomate.com/](https://make.powerautomate.com/)
2. Search for "Post to a channel when a webhook request is received" template
3. Select or create a connection to Microsoft Teams
4. Select **Create**
5. Choose your Team and Channel
6. Select **Create**
7. Select **Edit** from the menu
8. Select **When a Teams webhook request is received**
9. Copy the **HTTP POST URL** (this is your webhook URL)

**Method 3: Using Teams Workflows App**

1. Open **Workflows** application in Teams
2. Select **Create** at the top
3. Choose **Notifications** on the left
4. Select **Post to a channel when a webhook request is received**
5. Verify **Microsoft Teams** is enabled
6. Select **Next**
7. Choose your Team and Channel
8. Select **Create flow**
9. Copy the workflow URL
10. Select **Done**

**Security Note:** By default, anyone with the webhook URL can post messages. Treat your webhook URL as a secret!

**Valid Webhook URL Formats:**
- Power Automate: `https://prod*.apiflow.microsoft.com/...`
- O365 (deprecated): `https://outlook.office.com/webhook/...`
- O365 (deprecated): `https://*.webhook.office.com/...`

---

## 🏗️ Architecture

MultiChat MCP follows a clean, modular architecture designed for extensibility with **messenger-specific operations**:

```
┌─────────────────────────────────────────┐
│         MCP Client (Claude, etc)        │
└────────────────┬────────────────────────┘
                 │ MCP Protocol (stdio)
┌────────────────▼────────────────────────┐
│         MCP Server (server.go)          │
│   (Dynamically registers messenger      │
│    tools at runtime)                    │
└────────────────┬────────────────────────┘
                 │ Minimal Messenger Interface
┌────────────────▼────────────────────────┐
│  Platform Implementations (modular)     │
│  Each defines its OWN MCP operations    │
├──────────────────────────────────────────────────┤
│  ✅ WhatsApp │ ✅ Teams  │ ✅ Twitter/X │ 🔜 Telegram │
│  (7 tools)  │ (3 tools) │  (3 tools)   │  (8 tools)  │
└──────────────────────────────────────────────────┘
```

### Project Structure

```
multichatmcp/
├── main.go                      # Entry point & CLI
├── internal/
│   ├── messenger/
│   │   ├── interface.go         # Minimal messenger interface
│   │   ├── whatsapp/
│   │   │   ├── whatsapp.go      # WhatsApp implementation + MCP tools
│   │   │   └── types.go         # WhatsApp-specific types
│   │   ├── teams/
│   │   │   ├── teams.go         # Teams implementation + MCP tools
│   │   │   └── types.go         # Teams-specific types
│   │   └── twitter/
│   │       ├── twitter.go       # Twitter/X implementation + MCP tools
│   │       └── types.go         # Twitter/X-specific types
│   └── mcp/
│       └── server.go            # Generic MCP server
├── go.mod                       # Go dependencies
├── Makefile                     # Build automation
└── README.md                    # This file
```

### Key Design Principles

1. **Messenger-Specific Operations**: Each messenger platform registers its own unique set of MCP tools - there is **no common interface** for operations
2. **Minimal Interface**: The `Messenger` interface only requires `Connect()`, `Disconnect()`, `IsConnected()`, `GetMessengerName()`, and `RegisterMCPTools()`
3. **Dynamic Tool Registration**: MCP tools are registered at runtime based on the selected messenger type
4. **Platform Isolation**: Each messenger implementation is completely independent with its own types and operations
5. **Extensibility**: Adding new platforms means implementing new operations specific to that platform's capabilities

### Architecture Benefits

This **messenger-specific operations** architecture provides several key advantages:

#### ✅ True Platform Independence
Each messenger can expose operations that make sense for **that platform only**:
- WhatsApp has JID-based operations (`get_chat`, `get_direct_chat_by_contact`)
- Teams has webhook-based operations (`send_message`, `send_rich_message`, `validate_webhook`)
- Telegram could have poll/forward operations (`create_poll`, `forward_message`)

#### ✅ No Forced Abstractions
Platforms aren't forced to implement operations that don't make sense:
- A read-only analytics platform doesn't need `send_message`
- A simple notification service doesn't need complex chat management
- Each platform exposes exactly what it can do

#### ✅ Easy Evolution
Add new operations to specific platforms without affecting others:
- Add `create_poll` to Telegram without touching WhatsApp
- Implement `schedule_meeting` for Teams only
- Update one platform's types without breaking others

#### ✅ Runtime Flexibility
The MCP server automatically adapts to show only the operations available for the selected messenger:
```bash
# WhatsApp operations only
./multichat --messenger whatsapp

# Teams operations only
./multichat --messenger teams

# Each shows completely different MCP tools
```

---

## 🧩 Adding New Messaging Platforms

Want to add Teams, Telegram, Signal, or another platform? Each platform can define its own unique set of operations!

### Step 1: Define Platform-Specific Types

Create `internal/messenger/<platform>/types.go`:

```go
package platform

// Define your platform-specific types
type Contact struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    // Platform-specific fields
}

type Message struct {
    ID      string `json:"id"`
    Content string `json:"content"`
    // Platform-specific fields
}

// Add any platform-specific types you need
```

### Step 2: Implement the Minimal Interface

Create `internal/messenger/<platform>/<platform>.go`:

```go
package platform

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

type PlatformMessenger struct {
    // Your platform-specific client and state
}

func NewPlatformMessenger(config string) (*PlatformMessenger, error) {
    // Initialize your messenger
}

// Implement the minimal required methods:

func (p *PlatformMessenger) Connect(ctx context.Context) error {
    // Connect to your platform
}

func (p *PlatformMessenger) Disconnect() error {
    // Disconnect from your platform
}

func (p *PlatformMessenger) IsConnected() bool {
    // Return connection status
}

func (p *PlatformMessenger) GetMessengerName() string {
    return "platform-name" // e.g., "teams", "telegram"
}

// This is where you define YOUR platform's operations!
func (p *PlatformMessenger) RegisterMCPTools(mcpServer *server.MCPServer) {
    // Register ONLY the operations that make sense for your platform

    // Example: Teams might have channel operations
    mcpServer.AddTool(mcp.Tool{
        Name:        "list_channels",
        Description: "List all Teams channels",
        InputSchema: mcp.ToolInputSchema{
            Type: "object",
            Properties: map[string]interface{}{
                "team_id": map[string]interface{}{
                    "type":        "string",
                    "description": "The ID of the team",
                },
            },
            Required: []string{"team_id"},
        },
    }, p.handleListChannels)

    // Add more platform-specific tools...
}

// Implement your tool handlers
func (p *PlatformMessenger) handleListChannels(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Handle the operation
}
```

### Step 3: Register in Main

Add to the switch statement in `main.go`:

```go
case "teams":
    msg, err = teams.NewTeamsMessenger(deviceDB)
case "telegram":
    msg, err = telegram.NewTelegramMessenger(deviceDB)
```

### Step 4: Test & Document

- Test your platform-specific operations
- Update README with your platform's available MCP tools
- Document any platform-specific setup requirements
- Submit a PR! 🎉

### Example: Different Platforms, Different Operations

**WhatsApp** (7 operations):
- `search_contacts`, `list_messages`, `list_chats`, `get_chat`, `get_direct_chat_by_contact`, `get_contact_chats`, `send_message`

**Teams** (3 operations):
- `send_message`, `send_rich_message`, `validate_webhook`

**Twitter/X** (3 operations):
- `post_tweet`, `send_message`, `delete_tweet`

**Telegram** (hypothetical 8 operations):
- `list_chats`, `send_message`, `create_poll`, `pin_message`, `forward_message`, `get_chat_history`, `search_global`, `get_user_info`

Each platform implements **only what makes sense** for that platform!

---

## 🐛 Troubleshooting

<details>
<summary><strong>QR Code Not Appearing</strong></summary>

**Solution:**
```bash
# Run with debug logging
./multichat --messenger whatsapp --log-level debug

# Ensure no other WhatsApp Web sessions are active
# Try deleting device.db and restarting
rm device.db
./multichat --messenger whatsapp
```
</details>

<details>
<summary><strong>Connection Drops or Fails</strong></summary>

**Checklist:**
- ✅ Verify internet connectivity
- ✅ Check firewall/proxy settings
- ✅ Ensure WhatsApp is active on your phone
- ✅ Try reconnecting: delete `device.db` and re-scan QR code
</details>

<details>
<summary><strong>Database Locked Errors</strong></summary>

**Solution:**
```bash
# Ensure no other instances are running
pkill multichat

# Check file permissions
chmod 644 device.db

# If issue persists, remove and recreate
rm device.db
```
</details>

<details>
<summary><strong>MCP Client Can't Find Server</strong></summary>

**Solution:**
- Use **absolute paths** in configuration files
- Verify the binary is executable: `chmod +x /path/to/multichat`
- Check client logs for connection errors
- Test standalone: `./multichat --messenger whatsapp --log-level debug`
</details>

---

## ⚠️ Limitations & Known Issues

### 📜 Message History
WhatsApp's `whatsmeow` doesn't provide direct access to historical messages. To implement full history:

1. Set up event handlers for incoming messages
2. Store messages in a local database (SQLite, PostgreSQL, etc.)
3. Query from your database in `ListMessages`

**Example approach:**
```go
// Listen for new messages
client.AddEventHandler(func(evt interface{}) {
    if msg, ok := evt.(*events.Message); ok {
        // Store in your database
        db.SaveMessage(msg)
    }
})
```

### 🎬 Media Messages
Currently only **text messages** are supported. Media support (images, videos, documents) can be added by:
- Extending the `Message` struct
- Implementing download handlers
- Adding MCP tools for media retrieval

### 👥 Advanced Group Features
Not yet implemented:
- Create groups
- Add/remove members
- Update group settings
- Admin operations

These can be added by extending the `Messenger` interface and platform implementations.

---

## 🔨 Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

### Building

```bash
# Development build
make build

# Production build (optimized)
make build-prod

# Cross-compile for all platforms
make build-all

# Install to GOPATH/bin
make install

# Clean build artifacts
make clean
```

### Available Make Targets

```bash
make help           # Show all available targets
make build          # Build for current platform
make build-prod     # Build with optimizations (-ldflags="-s -w")
make build-all      # Cross-compile (Linux, macOS, Windows)
make install        # Install to GOPATH/bin
make clean          # Remove build artifacts
make test           # Run tests
```

---

## 🤝 Contributing

Contributions are **greatly appreciated**! Whether it's:

- 🐛 Bug reports
- 💡 Feature requests
- 📝 Documentation improvements
- 🔧 Code contributions
- 🌍 Platform integrations

### How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## 📜 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### What does this mean?

✅ Commercial use
✅ Modification
✅ Distribution
✅ Private use

---

## 🙏 Acknowledgments

This project stands on the shoulders of giants:

- **[whatsmeow](https://github.com/tulir/whatsmeow)** - Excellent WhatsApp Web multidevice library by [Tulir Asokan](https://github.com/tulir)
- **[go-teams-notify](https://github.com/atc0005/go-teams-notify)** - Microsoft Teams notification library by [atc0005](https://github.com/atc0005)
- **[gotwi](https://github.com/michimani/gotwi)** - Twitter API v2 library for Go by [michimani](https://github.com/michimani)
- **[mcp-go](https://github.com/modelcontextprotocol/go-mcp)** - Official Go implementation of MCP by [Mark3Labs](https://github.com/mark3labs)
- **[Cobra](https://github.com/spf13/cobra)** - Powerful CLI framework by [spf13](https://github.com/spf13)
- **[zerolog](https://github.com/rs/zerolog)** - Fast structured logger by [Olivier Poitrey](https://github.com/rs)

Special thanks to:
- The MCP community for the amazing protocol
- All contributors and testers
- You, for checking out this project! ⭐

---

## 📞 Support & Community

- 📫 **Issues**: [GitHub Issues](https://github.com/joao-costa/multichatmcp/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/joao-costa/multichatmcp/discussions)
- ⭐ **Star** this repo if you find it useful!

---

<div align="center">

**Made with ❤️ by [João Costa](https://github.com/joaocostamm)**

[⬆ Back to Top](#-multichat-mcp-server)

</div>
