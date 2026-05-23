// Package watcher polls configured dependencies at a set interval,
// fetches their changelogs, builds a digest, and dispatches notifications.
package watcher

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
	"github.com/yourorg/depwatch/internal/config"
	"github.com/yourorg/depwatch/internal/digest"
)

// Notifier is the interface satisfied by every notification back-end.
type Notifier interface {
	Send(subject, body string) error
}

// Fetcher retrieves raw changelog text for a given URL.
type Fetcher interface {
	Fetch(url string) (string, error)
}

// Watcher orchestrates the polling loop.
type Watcher struct {
	cfg      *config.Config
	fetcher  Fetcher
	parser   *changelog.Parser
	builder  *digest.Builder
	notifier Notifier
}

// New constructs a Watcher from its dependencies.
func New(cfg *config.Config, fetcher Fetcher, parser *changelog.Parser, builder *digest.Builder, notifier Notifier) *Watcher {
	return &Watcher{
		cfg:      cfg,
		fetcher:  fetcher,
		parser:   parser,
		builder:  builder,
		notifier: notifier,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(w.cfg.IntervalMinutes) * time.Minute)
	defer ticker.Stop()

	log.Printf("watcher: starting — interval %d min, %d dependencies",
		w.cfg.IntervalMinutes, len(w.cfg.Dependencies))

	// Run once immediately before waiting for the first tick.
	w.poll()

	for {
		select {
		case <-ticker.C:
			w.poll()
		case <-ctx.Done():
			log.Println("watcher: shutting down")
			return
		}
	}
}

// poll iterates over all configured dependencies and sends a digest when there
// are updates available.
func (w *Watcher) poll() {
	var updates []digest.Update

	for _, dep := range w.cfg.Dependencies {
		raw, err := w.fetcher.Fetch(dep.ChangelogURL)
		if err != nil {
			log.Printf("watcher: fetch %s: %v", dep.Name, err)
			continue
		}

		entries, err := w.parser.Parse(raw)
		if err != nil {
			log.Printf("watcher: parse %s: %v", dep.Name, err)
			continue
		}

		if len(entries) > 0 {
			updates = append(updates, digest.Update{Dependency: dep.Name, Entries: entries})
		}
	}

	if len(updates) == 0 {
		log.Println("watcher: no updates found")
		return
	}

	d := w.builder.Build(updates)
	if err := w.notifier.Send("depwatch digest", w.builder.FormatText(d)); err != nil {
		log.Printf("watcher: notify: %v", err)
	}
}
