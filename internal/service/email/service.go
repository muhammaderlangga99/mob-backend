package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Emailer defines SMTP-based delivery operations needed by the auth usecase.
type Emailer interface {
	SendVerificationEmail(to, token string) error
}

// SMTPEmailService handles sending emails via Gmail-compatible SMTP host.
type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
	appURL   string
}

// NewSMTPEmailService builds an SMTP client for verification emails.
func NewSMTPEmailService(host, port, username, password, from, appURL string) *SMTPEmailService {
	return &SMTPEmailService{host: host, port: port, username: username, password: password, from: from, appURL: strings.TrimRight(appURL, "/")}
}

// SendVerificationEmail dispatches the verification link to the given recipient.
func (s *SMTPEmailService) SendVerificationEmail(to, token string) error {
	if s.host == "" || s.port == "" || s.username == "" || s.password == "" || s.from == "" || s.appURL == "" {
		return fmt.Errorf("smtp configuration is incomplete")
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	link := fmt.Sprintf("%s/api/auth/verify-email?token=%s", s.appURL, token)
	msg := fmt.Sprintf("From: %s\r\n", s.from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += "Subject: Verifikasi Email Merchant Onboarding\r\n"
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/plain; charset=UTF-8\r\n\r\n"
	msg += fmt.Sprintf("Klik link ini untuk verifikasi: %s\r\n", link)

	return smtp.SendMail(fmt.Sprintf("%s:%s", s.host, s.port), auth, s.from, []string{to}, []byte(msg))
}
