package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSlackNotifier_EmptyURL(t *testing.T) {
	_, err := NewSlackNotifier("")
	if err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}
}

func TestNewSlackNotifier_ValidURL(t *testing.T) {
	n, err := NewSlackNotifier("https://hooks.slack.com/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSlackNotifier_Send_Success(t *testing.T) {
	var received string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		received = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewSlackNotifier(server.URL)
	if err := n.Send("hello depwatch"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = received
}

func TestSlackNotifier_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewSlackNotifier(server.URL)
	err := n.Send("test message")
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSlackNotifier_Send_InvalidURL(t *testing.T) {
	n, _ := NewSlackNotifier("http://127.0.0.1:0")
	err := n.Send("test")
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
