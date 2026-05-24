package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/depwatch/internal/changelog"
	"github.com/yourusername/depwatch/internal/config"
	"github.com/yourusername/depwatch/internal/digest"
	"github.com/yourusername/depwatch/internal/notifier"
	"github.com/yourusername/depwatch/internal/scheduler"
	"github.com/yourusername/depwatch/internal/watcher"
)

func main() {
	cfgPath := "depwatch.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fetcher := changelog.NewHTTPFetcher(0)
	parser := changelog.NewParser()
	builder := digest.NewBuilder()

	notify, err := buildNotifier(cfg)
	if err != nil {
		log.Fatalf("failed to create notifier: %v", err)
	}

	w := watcher.New(cfg, fetcher, parser, builder, notify)

	sched, err := scheduler.New(cfg.Interval, func(ctx context.Context) error {
		return w.Poll(ctx)
	})
	if err != nil {
		log.Fatalf("failed to create scheduler: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("shutting down depwatch...")
		cancel()
	}()

	log.Printf("depwatch started (interval: %s, dependencies: %d)", cfg.Interval, len(cfg.Dependencies))
	sched.Run(ctx)
}

// buildNotifier constructs the appropriate Notifier based on the configuration.
// It returns an error if the selected notifier cannot be initialised, and a
// fatal error if no notifier is configured at all.
func buildNotifier(cfg *config.Config) (notifier.Notifier, error) {
	if cfg.Slack.WebhookURL != "" {
		return notifier.NewSlackNotifier(cfg.Slack.WebhookURL)
	}
	if cfg.Email.Host != "" {
		return notifier.NewEmailNotifier(cfg.Email)
	}
	log.Fatal("no notifier configured: set slack.webhook_url or email.host")
	return nil, nil // unreachable, but satisfies the compiler
}
