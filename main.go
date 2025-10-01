package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/joao-costa/multichatmcp/internal/mcp"
	"github.com/joao-costa/multichatmcp/internal/messenger"
	"github.com/joao-costa/multichatmcp/internal/messenger/teams"
	"github.com/joao-costa/multichatmcp/internal/messenger/twitter"
	"github.com/joao-costa/multichatmcp/internal/messenger/whatsapp"
)

var (
	messengerType      string
	deviceDB           string
	webhookURL         string
	twitterAPIKey      string
	twitterAPISecret   string
	twitterToken       string
	twitterTokenSecret string
	logLevel           string
)

var rootCmd = &cobra.Command{
	Use:   "multichat",
	Short: "Multi-messenger MCP server",
	Long:  `A Model Context Protocol (MCP) server supporting multiple messaging platforms.`,
	RunE:  run,
}

func init() {
	rootCmd.Flags().StringVar(&messengerType, "messenger", "whatsapp", "Messenger type (whatsapp, teams, twitter)")
	rootCmd.Flags().StringVar(&deviceDB, "device", "device.db", "Device database file path (for WhatsApp)")
	rootCmd.Flags().StringVar(&webhookURL, "webhook", "", "Webhook URL (for Teams)")
	rootCmd.Flags().StringVar(&twitterAPIKey, "twitter-api-key", "", "Twitter API Key (for Twitter/X)")
	rootCmd.Flags().StringVar(&twitterAPISecret, "twitter-api-secret", "", "Twitter API Secret Key (for Twitter/X)")
	rootCmd.Flags().StringVar(&twitterToken, "twitter-token", "", "Twitter Access Token (for Twitter/X)")
	rootCmd.Flags().StringVar(&twitterTokenSecret, "twitter-token-secret", "", "Twitter Access Token Secret (for Twitter/X)")
	rootCmd.Flags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
}

func run(cmd *cobra.Command, args []string) error {
	// Setup logger
	setupLogger(logLevel)

	log.Info().
		Str("messenger", messengerType).
		Str("device_db", deviceDB).
		Str("log_level", logLevel).
		Msg("Starting MultiChat MCP Server")

	// Create messenger instance
	var msg messenger.Messenger
	var err error
	switch messengerType {
	case "whatsapp":
		msg, err = whatsapp.NewWhatsAppMessenger(deviceDB)
		if err != nil {
			return fmt.Errorf("failed to create WhatsApp messenger: %w", err)
		}
	case "teams":
		msg, err = teams.NewTeamsMessenger(webhookURL)
		if err != nil {
			return fmt.Errorf("failed to create Teams messenger: %w", err)
		}
	case "twitter":
		msg, err = twitter.NewTwitterMessenger(twitterAPIKey, twitterAPISecret, twitterToken, twitterTokenSecret)
		if err != nil {
			return fmt.Errorf("failed to create Twitter/X messenger: %w", err)
		}
	default:
		return fmt.Errorf("unsupported messenger type: %s", messengerType)
	}

	// Connect to messenger
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info().Msg("Connecting to messenger...")
	if err := msg.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to messenger: %w", err)
	}
	defer msg.Disconnect()

	log.Info().Msg("Messenger connected successfully")

	// Create and start MCP server
	mcpServer := mcp.NewServer(msg)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info().Msg("Shutdown signal received, closing connections...")
		cancel()
		msg.Disconnect()
		os.Exit(0)
	}()

	log.Info().Msg("Starting MCP server (stdio transport)...")
	if err := mcpServer.Serve(); err != nil {
		return fmt.Errorf("MCP server error: %w", err)
	}

	return nil
}

func setupLogger(level string) {
	// Set up zerolog with human-friendly output for development
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Set log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
