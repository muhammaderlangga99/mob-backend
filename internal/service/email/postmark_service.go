package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// PostmarkEmailService sends verification links via Postmark API.
type PostmarkEmailService struct {
	client      *http.Client
	serverToken string
	fromEmail   string
	appURL      string
	verifyPath  string
}

// NewPostmarkEmailService builds a Postmark emailer with configured env values.
func NewPostmarkEmailService(serverToken, fromEmail, appURL string) *PostmarkEmailService {
	return &PostmarkEmailService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		serverToken: strings.TrimSpace(serverToken),
		fromEmail:   strings.TrimSpace(fromEmail),
		appURL:      strings.TrimRight(appURL, "/"),
		verifyPath:  "/api/auth/verify-email",
	}
}

type postmarkPayload struct {
	From          string `json:"From"`
	To            string `json:"To"`
	Subject       string `json:"Subject"`
	HtmlBody      string `json:"HtmlBody"`
	MessageStream string `json:"MessageStream"`
}

// SendVerificationEmail constructs a Postmark payload and posts to the API.
func (s *PostmarkEmailService) SendVerificationEmail(toEmail, fullName, token string) error {
	if s.serverToken == "" || s.fromEmail == "" || s.appURL == "" {
		return fmt.Errorf("postmark configuration missing")
	}

	link := fmt.Sprintf("%s%s?token=%s", s.appURL, s.verifyPath, token)
	body := fmt.Sprintf("<p>Halo %s,</p><p>Klik tautan berikut untuk memverifikasi email Anda:</p><p><a href=\"%s\">%s</a></p><p>Terima kasih.</p>", fullName, link, link)

	payload := postmarkPayload{
		From:          s.fromEmail,
		To:            toEmail,
		Subject:       "Verifikasi Email Merchant Onboarding",
		HtmlBody:      body,
		MessageStream: "outbound",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.postmarkapp.com/email", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("X-Postmark-Server-Token", s.serverToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("postmark responded with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}
