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

	var notify notifier.Notifier
	if cfg.Slack.WebhookURL != "" {
		notify, err = notifier.NewSlackNotifier(cfg.Slack.WebhookURL)
		if err != nil {
			log.Fatalf("failed to create slack notifier: %v", err)
		}
	} else if cfg.Email.Host != "" {
		notify, err = notifier.NewEmailNotifier(cfg.Email)
		if err != nil {
			log.Fatalf("failed to create email notifier: %v", err)
		}
	} else {
		log.Fatal("no notifier configured: set slack.webhook_url or email.host")
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
