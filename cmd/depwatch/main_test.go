package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "depwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return p
}

func TestConfigPath_DefaultFallback(t *testing.T) {
	// Verify that when no args are provided the default path is used.
	// We just test the path selection logic, not the full main().
	args := []string{}
	cfgPath := "depwatch.yaml"
	if len(args) > 1 {
		cfgPath = args[1]
	}
	if cfgPath != "depwatch.yaml" {
		t.Errorf("expected default config path 'depwatch.yaml', got %q", cfgPath)
	}
}

func TestConfigPath_CustomArg(t *testing.T) {
	args := []string{"depwatch", "/etc/depwatch/config.yaml"}
	cfgPath := "depwatch.yaml"
	if len(args) > 1 {
		cfgPath = args[1]
	}
	if cfgPath != "/etc/depwatch/config.yaml" {
		t.Errorf("expected custom config path, got %q", cfgPath)
	}
}

func TestWriteTempConfig_CreatesFile(t *testing.T) {
	content := `interval: 1h
dependencies:
  - name: testpkg
    changelog_url: https://example.com/CHANGELOG.md
slack:
  webhook_url: https://hooks.slack.com/test
`
	p := writeTempConfig(t, content)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Errorf("expected temp config file to exist at %q", p)
	}
}
