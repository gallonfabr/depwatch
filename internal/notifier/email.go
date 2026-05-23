package notifier

import (
	"errors"
	"fmt"
	"net/smtp"
)

// EmailNotifier sends digest notifications via SMTP email.
type EmailNotifier struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       []string
}

// NewEmailNotifier creates a new EmailNotifier. Returns an error if required
// fields are missing.
func NewEmailNotifier(host string, port int, username, password, from string, to []string) (*EmailNotifier, error) {
	if host == "" {
		return nil, errors.New("email notifier: host must not be empty")
	}
	if from == "" {
		return nil, errors.New("email notifier: from address must not be empty")
	}
	if len(to) == 0 {
		return nil, errors.New("email notifier: at least one recipient required")
	}
	return &EmailNotifier{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}, nil
}

// Send delivers the digest text as an email to all configured recipients.
func (e *EmailNotifier) Send(subject, body string) error {
	addr := fmt.Sprintf("%s:%d", e.host, e.port)

	var auth smtp.Auth
	if e.username != "" {
		auth = smtp.PlainAuth("", e.username, e.password, e.host)
	}

	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n",
		e.from, e.to[0], subject,
	)
	message := []byte(header + body)

	if err := smtp.SendMail(addr, auth, e.from, e.to, message); err != nil {
		return fmt.Errorf("email notifier: send failed: %w", err)
	}
	return nil
}
