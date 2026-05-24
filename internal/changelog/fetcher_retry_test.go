package changelog

import (
	"errors"
	"testing"
	"time"
)

type stubFetcher struct {
	calls    int
	failUntil int
	response string
	err      error
}

func (s *stubFetcher) Fetch(_ string) (string, error) {
	s.calls++
	if s.calls <= s.failUntil {
		return "", s.err
	}
	return s.response, nil
}

func TestNewRetryFetcher_NilInner(t *testing.T) {
	_, err := NewRetryFetcher(nil, DefaultRetryConfig())
	if err == nil {
		t.Fatal("expected error for nil inner fetcher")
	}
}

func TestNewRetryFetcher_Valid(t *testing.T) {
	stub := &stubFetcher{}
	rf, err := NewRetryFetcher(stub, DefaultRetryConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rf == nil {
		t.Fatal("expected non-nil RetryFetcher")
	}
}

func TestRetryFetcher_Fetch_SuccessFirstTry(t *testing.T) {
	stub := &stubFetcher{response: "changelog content"}
	rf, _ := NewRetryFetcher(stub, RetryConfig{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond})
	got, err := rf.Fetch("http://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "changelog content" {
		t.Fatalf("unexpected response: %q", got)
	}
	if stub.calls != 1 {
		t.Fatalf("expected 1 call, got %d", stub.calls)
	}
}

func TestRetryFetcher_Fetch_RetriesAndSucceeds(t *testing.T) {
	stub := &stubFetcher{failUntil: 2, response: "ok", err: errors.New("timeout")}
	rf, _ := NewRetryFetcher(stub, RetryConfig{MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond})
	got, err := rf.Fetch("http://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ok" {
		t.Fatalf("unexpected response: %q", got)
	}
	if stub.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", stub.calls)
	}
}

func TestRetryFetcher_Fetch_AllAttemptsFail(t *testing.T) {
	fetchErr := errors.New("connection refused")
	stub := &stubFetcher{failUntil: 10, err: fetchErr}
	rf, _ := NewRetryFetcher(stub, RetryConfig{MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond})
	_, err := rf.Fetch("http://example.com")
	if !errors.Is(err, fetchErr) {
		t.Fatalf("expected fetchErr, got %v", err)
	}
	if stub.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", stub.calls)
	}
}
