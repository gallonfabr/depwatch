package watcher_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/yourorg/depwatch/internal/store"
	"github.com/yourorg/depwatch/internal/watcher"
)

// stubFetcherOnce returns content on the first call and an error on subsequent
// calls, letting us verify that the watcher skips already-seen versions.
type stubFetcherOnce struct {
	content string
	calls   int
}

func (f *stubFetcherOnce) Fetch(_ string) (string, error) {
	f.calls++
	if f.calls == 1 {
		return f.content, nil
	}
	return "", errors.New("should not be called again")
}

func TestWatcher_UsesStore_SkipsSeen(t *testing.T) {
	st, err := store.New(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}

	cfg := sampleConfig()
	fetcher := &stubFetcherOnce{content: "## [1.2.3] - 2024-01-01\n\n- fix: something"}
	notified := 0
	notifier := &mockNotifier{sendFn: func(_ string) error { notified++; return nil }}

	w := watcher.New(cfg, fetcher, notifier, watcher.WithStore(st))

	// First poll — new version, should notify.
	if err := w.Poll(); err != nil {
		t.Fatalf("first Poll error: %v", err)
	}
	if notified != 1 {
		t.Errorf("expected 1 notification after first poll, got %d", notified)
	}

	// Second poll — same version already stored, should NOT notify.
	if err := w.Poll(); err != nil {
		t.Fatalf("second Poll error: %v", err)
	}
	if notified != 1 {
		t.Errorf("expected still 1 notification after second poll, got %d", notified)
	}
}
