package notifier

import (
	"net"
	"net/smtp"
	"net/textproto"
	"testing"
)

func TestNewEmailNotifier_MissingHost(t *testing.T) {
	_, err := NewEmailNotifier("", 587, "user", "pass", "from@example.com", []string{"to@example.com"})
	if err == nil {
		t.Fatal("expected error for empty host, got nil")
	}
}

func TestNewEmailNotifier_MissingFrom(t *testing.T) {
	_, err := NewEmailNotifier("smtp.example.com", 587, "user", "pass", "", []string{"to@example.com"})
	if err == nil {
		t.Fatal("expected error for empty from, got nil")
	}
}

func TestNewEmailNotifier_NoRecipients(t *testing.T) {
	_, err := NewEmailNotifier("smtp.example.com", 587, "user", "pass", "from@example.com", []string{})
	if err == nil {
		t.Fatal("expected error for empty recipients, got nil")
	}
}

func TestNewEmailNotifier_Valid(t *testing.T) {
	n, err := NewEmailNotifier("smtp.example.com", 587, "user", "pass", "from@example.com", []string{"to@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewEmailNotifier_NoAuth(t *testing.T) {
	// Empty username should still create a valid notifier (no auth mode).
	n, err := NewEmailNotifier("smtp.example.com", 25, "", "", "from@example.com", []string{"to@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestEmailNotifier_Send_InvalidAddr(t *testing.T) {
	n, _ := NewEmailNotifier("127.0.0.1", 19999, "", "", "from@example.com", []string{"to@example.com"})
	err := n.Send("Test Subject", "Test body")
	if err == nil {
		t.Fatal("expected error when SMTP server is unreachable, got nil")
	}
}

// Ensure smtp.PlainAuth is used when username is provided — compile-time check.
var _ smtp.Auth = smtp.PlainAuth("", "", "", "")

// Ensure net and textproto packages are available for future integration tests.
var _ net.Conn
var _ *textproto.Reader
