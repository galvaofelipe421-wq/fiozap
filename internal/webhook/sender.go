package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	httpTimeout     = 10 * time.Second
	contentTypeJSON = "application/json"
	userAgent       = "FioZap-Webhook/1.0"
)

type Sender struct {
	client *http.Client
}

func NewSender() *Sender {
	return &Sender{
		client: &http.Client{Timeout: httpTimeout},
	}
}

type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func (s *Sender) Send(ctx context.Context, url string, payload *WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("User-Agent", userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
