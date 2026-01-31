package email

// Emailer defines delivery operations needed by auth usecase.
type Emailer interface {
	SendVerificationEmail(toEmail, fullName, token string) error
}
