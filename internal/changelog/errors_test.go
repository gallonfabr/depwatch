package changelog_test

import (
	"testing"

	"github.com/yourorg/depwatch/internal/changelog"
)

func TestSentinelErrors_NotNil(t *testing.T) {
	errors := []error{
		changelog.ErrMissingURL,
		changelog.ErrMissingOwner,
		changelog.ErrMissingRepo,
		changelog.ErrUnknownSourceType,
	}
	for _, err := range errors {
		if err == nil {
			t.Errorf("expected non-nil sentinel error, got nil")
		}
	}
}

func TestSentinelErrors_UniqueMessages(t *testing.T) {
	seen := map[string]bool{}
	for _, err := range []error{
		changelog.ErrMissingURL,
		changelog.ErrMissingOwner,
		changelog.ErrMissingRepo,
		changelog.ErrUnknownSourceType,
	} {
		msg := err.Error()
		if seen[msg] {
			t.Errorf("duplicate error message: %q", msg)
		}
		seen[msg] = true
	}
}
