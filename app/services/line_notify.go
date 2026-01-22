package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"timeLedger/app"
)

type LINENotifyService interface {
	SendMessage(ctx context.Context, token string, message string) error
	SendSticker(ctx context.Context, token, message, stickerPackageID, stickerID string) error
}

type LINENotifyServiceImpl struct {
	app    *app.App
	apiURL string
	client *http.Client
}

type LINESticker struct {
	PackageID string `json:"packageId"`
	ID        string `json:"stickerId"`
}

type LINENotifyRequest struct {
	Message string       `json:"message"`
	Sticker *LINESticker `json:"sticker,omitempty"`
}

type LINENotifyResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewLINENotifyService(app *app.App) LINENotifyService {
	return &LINENotifyServiceImpl{
		app:    app,
		apiURL: "https://notify-api.line.me/api/notify",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *LINENotifyServiceImpl) SendMessage(ctx context.Context, token string, message string) error {
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	reqBody := LINENotifyRequest{
		Message: message,
	}

	return s.sendRequest(ctx, token, reqBody)
}

func (s *LINENotifyServiceImpl) SendSticker(ctx context.Context, token string, message string, stickerPackageID string, stickerID string) error {
	if stickerPackageID == "" || stickerID == "" {
		return fmt.Errorf("sticker package ID and sticker ID cannot be empty")
	}

	reqBody := LINENotifyRequest{
		Message: message,
		Sticker: &LINESticker{
			PackageID: stickerPackageID,
			ID:        stickerID,
		},
	}

	return s.sendRequest(ctx, token, reqBody)
}

func (s *LINENotifyServiceImpl) sendRequest(ctx context.Context, token string, reqBody LINENotifyRequest) error {
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var lineResp LINENotifyResponse
		if err := json.Unmarshal(body, &lineResp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return fmt.Errorf("LINE Notify API error (status %d): %s", lineResp.Status, lineResp.Message)
	}

	return nil
}
