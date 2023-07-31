package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Webhook is the slack post message api via webhook.
type Webhook struct {
	url     string
	channel string
}

// NewWebhook creates the webhook with the config.
func NewWebhook(config *Config) *Webhook {
	return &Webhook{
		url:     config.WebhookURL,
		channel: config.Channel,
	}
}

// PostMessage sends a message to the Slack channel using the configured
// webhook.
func (w *Webhook) PostMessage(msg string) error {
	payload := map[string]string{
		"text":    msg,
		"channel": w.channel,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		w.url,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful (status code 200-299).
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"failed to post message to Slack. Status code: %d",
			resp.StatusCode,
		)
	}

	return nil
}
