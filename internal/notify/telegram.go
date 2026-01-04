package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const telegramAPIURL = "https://api.telegram.org"

// TelegramConfig holds Telegram bot configuration.
type TelegramConfig struct {
	BotToken string
	ChatID   string
}

// Telegram implements Notifier for Telegram Bot API.
type Telegram struct {
	config     TelegramConfig
	httpClient *http.Client
}

// NewTelegram creates a new Telegram notifier.
func NewTelegram(cfg TelegramConfig) *Telegram {
	return &Telegram{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// sendMessageRequest is the request body for Telegram sendMessage API.
type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// telegramResponse is the generic Telegram API response.
type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

// Send sends a message to the configured Telegram chat.
func (t *Telegram) Send(ctx context.Context, message string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramAPIURL, t.config.BotToken)

	reqBody := sendMessageRequest{
		ChatID:    t.config.ChatID,
		Text:      message,
		ParseMode: "HTML",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var tgResp telegramResponse
	if err := json.Unmarshal(body, &tgResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !tgResp.OK {
		return fmt.Errorf("telegram API error: %s", tgResp.Description)
	}

	return nil
}

// Validate checks that the Telegram configuration is valid.
func (c TelegramConfig) Validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if c.ChatID == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID is required")
	}
	return nil
}
