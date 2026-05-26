package changelog

import (
	"testing"
	"time"
)

var sampleRedactorEntries = []Entry{
	{
		Dependency: "libfoo",
		Version:    "1.0.0",
		Date:       time.Now(),
		Body:       "Fixed bug. Token: ghp_abc123XYZ secret",
		Link:       "https://internal.corp.example.com/releases/1.0.0",
	},
	{
		Dependency: "libbar",
		Version:    "2.1.0",
		Date:       time.Now(),
		Body:       "New feature added.",
		Link:       "https://github.com/libbar/releases",
	},
}

func TestNewRedactor_NotNil(t *testing.T) {
	r := NewRedactor()
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestRedactor_NoPatterns_ReturnsUnchanged(t *testing.T) {
	r := NewRedactor()
	out := r.Apply(sampleRedactorEntries)
	if out[0].Body != sampleRedactorEntries[0].Body {
		t.Errorf("expected body unchanged, got %q", out[0].Body)
	}
}

func TestRedactor_TokenPattern_RedactsBody(t *testing.T) {
	r := NewRedactor(WithRedactPattern(`ghp_[A-Za-z0-9]+`))
	out := r.Apply(sampleRedactorEntries)
	if got := out[0].Body; got == sampleRedactorEntries[0].Body {
		t.Error("expected body to be redacted")
	}
	const want = "[REDACTED]"
	for _, part := range []string{"ghp_abc123XYZ"} {
		if contains := func(s, sub string) bool {
			return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
		}; contains(out[0].Body, part) {
			t.Errorf("sensitive token %q still present in body", part)
		}
		_ = want
	}
}

func TestRedactor_CustomPlaceholder(t *testing.T) {
	r := NewRedactor(
		WithRedactPattern(`ghp_[A-Za-z0-9]+`),
		WithRedactPlaceholder("***"),
	)
	out := r.Apply(sampleRedactorEntries)
	if out[0].Body == sampleRedactorEntries[0].Body {
		t.Error("expected body to change")
	}
}

func TestRedactor_LinkRedacted(t *testing.T) {
	r := NewRedactor(WithRedactPattern(`internal\.corp\.example\.com`))
	out := r.Apply(sampleRedactorEntries)
	if containsSubstring(out[0].Link, "internal.corp.example.com") {
		t.Error("expected link to be redacted")
	}
}

func TestRedactor_InvalidPattern_Ignored(t *testing.T) {
	r := NewRedactor(WithRedactPattern(`[invalid`))
	if len(r.patterns) != 0 {
		t.Error("expected invalid pattern to be ignored")
	}
}

func TestRedactor_UnaffectedEntry_Unchanged(t *testing.T) {
	r := NewRedactor(WithRedactPattern(`ghp_[A-Za-z0-9]+`))
	out := r.Apply(sampleRedactorEntries)
	if out[1].Body != sampleRedactorEntries[1].Body {
		t.Errorf("unaffected entry body changed: %q", out[1].Body)
	}
}

func containsSubstring(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) &&
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}()
}
