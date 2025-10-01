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

MultiChat MCP Server is a powerful Go-based implementation of the [Model Context Protocol (MCP)](https://modelcontextprotocol.io) that enables AI assistants like Claude, GPT, and others to seamlessly interact with your messaging platforms. Start with WhatsApp and expand to Telegram, Signal, and more with our modular architecture.

### Why MultiChat MCP?

- 🤖 **AI-Native**: Built specifically for AI assistants to read and send messages
- 🔌 **Plug & Play**: Easy integration with Claude Desktop, Cursor, and any MCP-compatible client
- 🧩 **Modular Design**: Clean interface makes adding new platforms straightforward
- 🔒 **Privacy First**: Your data stays on your machine
- ⚡ **Lightning Fast**: Written in Go for optimal performance

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 📱 Platform Support
- ✅ **WhatsApp** (via [whatsmeow](https://github.com/tulir/whatsmeow))
- 🔜 Telegram (coming soon)
- 🔜 Signal (coming soon)
- 🔜 Discord (planned)

</td>
<td width="50%">

### 🛠️ Core Capabilities
- 💬 Send & receive messages
- 👥 Contact & chat management
- 🔍 Full-text message search
- 📊 Pagination & filtering
- 🕐 Time-based queries

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

### First Run - WhatsApp Setup

```bash
./multichat --messenger whatsapp --device mydevice.db --log-level debug
```

**🔐 Authentication Steps:**
1. A QR code will appear in your terminal
2. Open WhatsApp on your phone → **Settings** → **Linked Devices** → **Link a Device**
3. Scan the QR code
4. Done! Your session is saved for future use

---

## 📖 Documentation

### Command-Line Usage

```bash
./multichat [flags]

Flags:
  --messenger string    Messaging platform to use (default "whatsapp")
  --device string       Device database file path (default "device.db")
  --log-level string    Logging level: debug, info, warn, error (default "info")
  -h, --help           Show help information
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

## 🏗️ Architecture

MultiChat MCP follows a clean, modular architecture designed for extensibility:

```
┌─────────────────────────────────────────┐
│         MCP Client (Claude, etc)        │
└────────────────┬────────────────────────┘
                 │ MCP Protocol (stdio)
┌────────────────▼────────────────────────┐
│         MCP Server (server.go)          │
└────────────────┬────────────────────────┘
                 │ Messenger Interface
┌────────────────▼────────────────────────┐
│  Platform Implementations (modular)     │
├─────────────────────────────────────────┤
│  ✅ WhatsApp  │  🔜 Telegram  │  🔜 Signal │
└─────────────────────────────────────────┘
```

### Project Structure

```
multichatmcp/
├── main.go                      # Entry point & CLI
├── internal/
│   ├── messenger/
│   │   ├── interface.go         # Messenger interface definition
│   │   └── whatsapp/
│   │       └── whatsapp.go      # WhatsApp implementation
│   └── mcp/
│       └── server.go            # MCP server implementation
├── go.mod                       # Go dependencies
├── Makefile                     # Build automation
└── README.md                    # This file
```

### Key Design Principles

1. **Interface-Driven**: All platforms implement a common `Messenger` interface
2. **Dependency Injection**: Platform implementations are injected into MCP server
3. **Separation of Concerns**: CLI, MCP layer, and platform logic are isolated
4. **Extensibility**: Adding new platforms requires minimal changes

---

## 🧩 Adding New Messaging Platforms

Want to add Telegram, Signal, or another platform? Here's how:

### Step 1: Implement the Interface

Create `internal/messenger/<platform>/<platform>.go`:

```go
package platform

import (
    "context"
    "github.com/joao-costa/multichatmcp/internal/messenger"
)

type PlatformMessenger struct {
    // Your implementation
}

func NewPlatformMessenger(config string) (messenger.Messenger, error) {
    // Initialize your messenger
}

// Implement all messenger.Messenger interface methods:
// - Connect(ctx context.Context) error
// - Disconnect() error
// - SearchContacts(ctx context.Context, query string) ([]Contact, error)
// - ListMessages(ctx context.Context, filter MessageFilter) ([]Message, error)
// - ListChats(ctx context.Context, limit, page int) ([]Chat, error)
// - GetChat(ctx context.Context, chatJID string) (*Chat, error)
// - GetDirectChatByContact(ctx context.Context, phoneNumber string) (*Chat, error)
// - GetContactChats(ctx context.Context, contactJID string) ([]Chat, error)
// - SendMessage(ctx context.Context, recipient, message string) error
// - IsConnected() bool
```

### Step 2: Register in Main

Add to the switch statement in `main.go`:

```go
case "telegram":
    msg, err = telegram.NewTelegramMessenger(deviceDB)
case "signal":
    msg, err = signal.NewSignalMessenger(deviceDB)
```

### Step 3: Test & Document

- Add tests for your implementation
- Update README with platform-specific setup
- Submit a PR! 🎉

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
