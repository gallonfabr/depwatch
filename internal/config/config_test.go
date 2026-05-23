package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "depwatch.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
check_interval: 12h
dependencies:
  - name: react
    ecosystem: npm
    version: "18.0.0"
  - name: requests
    ecosystem: pypi
    version: "2.28.0"
notifiers:
  slack:
    webhook_url: https://hooks.slack.com/xxx
    channel: "#deps"
`
	path := writeTemp(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CheckInterval != 12*time.Hour {
		t.Errorf("expected 12h interval, got %v", cfg.CheckInterval)
	}
	if len(cfg.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(cfg.Dependencies))
	}
	if cfg.Notifiers.Slack.Channel != "#deps" {
		t.Errorf("unexpected slack channel: %q", cfg.Notifiers.Slack.Channel)
	}
}

func TestLoad_DefaultInterval(t *testing.T) {
	content := `
dependencies:
  - name: lodash
    ecosystem: npm
`
	path := writeTemp(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CheckInterval != 24*time.Hour {
		t.Errorf("expected default 24h interval, got %v", cfg.CheckInterval)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/depwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoDependencies(t *testing.T) {
	content := `check_interval: 1h\n`
	path := writeTemp(t, content)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty dependencies")
	}
}

func TestLoad_MissingEcosystem(t *testing.T) {
	content := `
dependencies:
  - name: react
`
	path := writeTemp(t, content)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing ecosystem")
	}
}
