package watcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
	"github.com/yourorg/depwatch/internal/config"
	"github.com/yourorg/depwatch/internal/digest"
	"github.com/yourorg/depwatch/internal/watcher"
)

// --- stubs ---

type stubFetcher struct{ content string; err error }

func (s *stubFetcher) Fetch(_ string) (string, error) { return s.content, s.err }

type stubNotifier struct{ calls int; lastBody string; err error }

func (s *stubNotifier) Send(_, body string) error {
	s.calls++
	s.lastBody = body
	return s.err
}

func sampleConfig(intervalMin int) *config.Config {
	return &config.Config{
		IntervalMinutes: intervalMin,
		Dependencies: []config.Dependency{
			{Name: "mylib", ChangelogURL: "https://example.com/CHANGELOG.md"},
		},
	}
}

const sampleChangelog = `# Changelog

## [1.2.0] - 2024-05-01
### Added
- new feature

## [1.1.0] - 2024-04-01
### Fixed
- bug fix
`

// --- tests ---

func TestWatcher_Poll_SendsNotification(t *testing.T) {
	notifier := &stubNotifier{}
	w := watcher.New(
		sampleConfig(60),
		&stubFetcher{content: sampleChangelog},
		changelog.NewParser(),
		digest.NewBuilder(),
		notifier,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if notifier.calls == 0 {
		t.Fatal("expected at least one notification to be sent")
	}
	if notifier.lastBody == "" {
		t.Error("expected non-empty notification body")
	}
}

func TestWatcher_Poll_FetchError_NoNotification(t *testing.T) {
	notifier := &stubNotifier{}
	w := watcher.New(
		sampleConfig(60),
		&stubFetcher{err: errors.New("network error")},
		changelog.NewParser(),
		digest.NewBuilder(),
		notifier,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if notifier.calls != 0 {
		t.Fatalf("expected 0 notifications on fetch error, got %d", notifier.calls)
	}
}

func TestWatcher_Poll_EmptyChangelog_NoNotification(t *testing.T) {
	notifier := &stubNotifier{}
	w := watcher.New(
		sampleConfig(60),
		&stubFetcher{content: ""},
		changelog.NewParser(),
		digest.NewBuilder(),
		notifier,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if notifier.calls != 0 {
		t.Fatalf("expected 0 notifications for empty changelog, got %d", notifier.calls)
	}
}
