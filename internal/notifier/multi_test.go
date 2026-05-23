package notifier

import (
	"errors"
	"testing"
)

// stubNotifier is a test double that records calls and optionally returns an error.
type stubNotifier struct {
	called  bool
	subject string
	body    string
	errOut  error
}

func (s *stubNotifier) Send(subject, body string) error {
	s.called = true
	s.subject = subject
	s.body = body
	return s.errOut
}

func TestNewMultiNotifier_Empty(t *testing.T) {
	_, err := NewMultiNotifier()
	if err == nil {
		t.Fatal("expected error for empty notifier list, got nil")
	}
}

func TestNewMultiNotifier_Valid(t *testing.T) {
	n := &stubNotifier{}
	mn, err := NewMultiNotifier(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mn == nil {
		t.Fatal("expected non-nil MultiNotifier")
	}
}

func TestMultiNotifier_Send_AllSucceed(t *testing.T) {
	a, b := &stubNotifier{}, &stubNotifier{}
	mn, _ := NewMultiNotifier(a, b)

	if err := mn.Send("subj", "body"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("expected both notifiers to be called")
	}
	if a.subject != "subj" || a.body != "body" {
		t.Errorf("unexpected payload on notifier a: %q / %q", a.subject, a.body)
	}
}

func TestMultiNotifier_Send_PartialFailure(t *testing.T) {
	a := &stubNotifier{errOut: errors.New("slack down")}
	b := &stubNotifier{}
	mn, _ := NewMultiNotifier(a, b)

	err := mn.Send("subj", "body")
	if err == nil {
		t.Fatal("expected error when one notifier fails")
	}
	if !b.called {
		t.Error("expected second notifier to still be called after first failure")
	}
}

func TestMultiNotifier_Send_AllFail(t *testing.T) {
	a := &stubNotifier{errOut: errors.New("err a")}
	b := &stubNotifier{errOut: errors.New("err b")}
	mn, _ := NewMultiNotifier(a, b)

	err := mn.Send("subj", "body")
	if err == nil {
		t.Fatal("expected error when all notifiers fail")
	}
}
