package twitter

import "time"

// TwitterConfig holds configuration for Twitter/X messenger
type TwitterConfig struct {
	APIKey            string `json:"api_key"`
	APISecretKey      string `json:"api_secret_key"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

// Tweet represents a Twitter/X tweet
type Tweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TweetResponse represents the response from sending a tweet
type TweetResponse struct {
	TweetID string `json:"tweet_id"`
	Text    string `json:"text"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// User represents a Twitter/X user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// DirectMessage represents a Twitter/X DM
type DirectMessage struct {
	ID             string    `json:"id"`
	Text           string    `json:"text"`
	SenderID       string    `json:"sender_id"`
	RecipientID    string    `json:"recipient_id"`
	CreatedAt      time.Time `json:"created_at"`
	ConversationID string    `json:"conversation_id"`
}

// DMResponse represents the response from sending a DM
type DMResponse struct {
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}
