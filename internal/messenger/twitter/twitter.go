package twitter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"github.com/rs/zerolog/log"
)

// TwitterMessenger implements the Messenger interface for Twitter/X
type TwitterMessenger struct {
	config    TwitterConfig
	client    *gotwi.Client
	connected bool
}

// NewTwitterMessenger creates a new Twitter/X messenger instance
func NewTwitterMessenger(apiKey, apiSecretKey, accessToken, accessTokenSecret string) (*TwitterMessenger, error) {
	if apiKey == "" || apiSecretKey == "" || accessToken == "" || accessTokenSecret == "" {
		return nil, fmt.Errorf("all Twitter API credentials are required")
	}

	return &TwitterMessenger{
		config: TwitterConfig{
			APIKey:            apiKey,
			APISecretKey:      apiSecretKey,
			AccessToken:       accessToken,
			AccessTokenSecret: accessTokenSecret,
		},
		connected: false,
	}, nil
}

// Connect establishes connection to Twitter/X API
func (t *TwitterMessenger) Connect(ctx context.Context) error {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		APIKey:               t.config.APIKey,
		APIKeySecret:         t.config.APISecretKey,
		OAuthToken:           t.config.AccessToken,
		OAuthTokenSecret:     t.config.AccessTokenSecret,
	}

	client, err := gotwi.NewClient(in)
	if err != nil {
		return fmt.Errorf("failed to create Twitter client: %w", err)
	}

	t.client = client
	t.connected = true
	log.Info().Msg("Twitter/X messenger connected successfully")
	return nil
}

// Disconnect closes the Twitter/X connection
func (t *TwitterMessenger) Disconnect() error {
	t.connected = false
	t.client = nil
	log.Info().Msg("Twitter/X messenger disconnected")
	return nil
}

// IsConnected returns the connection status
func (t *TwitterMessenger) IsConnected() bool {
	return t.connected && t.client != nil
}

// GetMessengerName returns the name of the messenger platform
func (t *TwitterMessenger) GetMessengerName() string {
	return "twitter"
}

// postTweet posts a tweet to Twitter/X
func (t *TwitterMessenger) postTweet(ctx context.Context, text string, replyToTweetID string) (*TweetResponse, error) {
	if !t.IsConnected() {
		return nil, fmt.Errorf("not connected to Twitter/X")
	}

	if text == "" {
		return nil, fmt.Errorf("tweet text cannot be empty")
	}

	// Twitter has a 280 character limit
	if len([]rune(text)) > 280 {
		return nil, fmt.Errorf("tweet text exceeds 280 characters")
	}

	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	// Add reply reference if provided
	if replyToTweetID != "" {
		p.Reply = &types.CreateInputReply{
			InReplyToTweetID: replyToTweetID,
		}
	}

	res, err := managetweet.Create(ctx, t.client, p)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post tweet")
		return &TweetResponse{
			Text:    text,
			Success: false,
			Error:   err.Error(),
		}, err
	}

	log.Info().Str("tweet_id", gotwi.StringValue(res.Data.ID)).Msg("Tweet posted successfully")
	return &TweetResponse{
		TweetID: gotwi.StringValue(res.Data.ID),
		Text:    gotwi.StringValue(res.Data.Text),
		Success: true,
	}, nil
}

// deleteTweet deletes a tweet by ID
func (t *TwitterMessenger) deleteTweet(ctx context.Context, tweetID string) error {
	if !t.IsConnected() {
		return fmt.Errorf("not connected to Twitter/X")
	}

	if tweetID == "" {
		return fmt.Errorf("tweet ID cannot be empty")
	}

	p := &types.DeleteInput{
		ID: tweetID,
	}

	_, err := managetweet.Delete(ctx, t.client, p)
	if err != nil {
		log.Error().Err(err).Str("tweet_id", tweetID).Msg("Failed to delete tweet")
		return fmt.Errorf("failed to delete tweet: %w", err)
	}

	log.Info().Str("tweet_id", tweetID).Msg("Tweet deleted successfully")
	return nil
}

// RegisterMCPTools registers Twitter/X-specific MCP tools
func (t *TwitterMessenger) RegisterMCPTools(mcpServer *server.MCPServer) {
	// post_tweet - Post a tweet
	mcpServer.AddTool(mcp.Tool{
		Name:        "post_tweet",
		Description: "Post a tweet to Twitter/X (max 280 characters)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"text": map[string]interface{}{
					"type":        "string",
					"description": "The text content of the tweet (max 280 characters)",
				},
				"reply_to_tweet_id": map[string]interface{}{
					"type":        "string",
					"description": "Optional: ID of the tweet to reply to",
				},
			},
			Required: []string{"text"},
		},
	}, t.handlePostTweet)

	// send_message - Alias for post_tweet to match other messengers
	mcpServer.AddTool(mcp.Tool{
		Name:        "send_message",
		Description: "Send a message (tweet) to Twitter/X (max 280 characters)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "The text content of the tweet (max 280 characters)",
				},
				"reply_to_tweet_id": map[string]interface{}{
					"type":        "string",
					"description": "Optional: ID of the tweet to reply to",
				},
			},
			Required: []string{"message"},
		},
	}, t.handleSendMessage)

	// delete_tweet - Delete a tweet
	mcpServer.AddTool(mcp.Tool{
		Name:        "delete_tweet",
		Description: "Delete a tweet by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"tweet_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the tweet to delete",
				},
			},
			Required: []string{"tweet_id"},
		},
	}, t.handleDeleteTweet)
}

// Tool handlers

func (t *TwitterMessenger) handlePostTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Text           string `json:"text"`
		ReplyToTweetID string `json:"reply_to_tweet_id"`
	}

	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	result, err := t.postTweet(ctx, args.Text, args.ReplyToTweetID)
	if err != nil {
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultError(fmt.Sprintf("failed to post tweet: %v\nDetails: %s", err, string(resultJSON))), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (t *TwitterMessenger) handleSendMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Message        string `json:"message"`
		ReplyToTweetID string `json:"reply_to_tweet_id"`
	}

	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	result, err := t.postTweet(ctx, args.Message, args.ReplyToTweetID)
	if err != nil {
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultError(fmt.Sprintf("failed to send message: %v\nDetails: %s", err, string(resultJSON))), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (t *TwitterMessenger) handleDeleteTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		TweetID string `json:"tweet_id"`
	}

	argsBytes, _ := json.Marshal(request.Params.Arguments)
	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	err := t.deleteTweet(ctx, args.TweetID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete tweet: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"success": true, "tweet_id": "%s", "message": "Tweet deleted successfully"}`, args.TweetID)), nil
}
