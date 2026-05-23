package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level depwatch configuration.
type Config struct {
	CheckInterval time.Duration `yaml:"check_interval"`
	Dependencies  []Dependency  `yaml:"dependencies"`
	Notifiers     Notifiers     `yaml:"notifiers"`
}

// Dependency describes a single dependency to watch.
type Dependency struct {
	Name     string `yaml:"name"`
	Ecosystem string `yaml:"ecosystem"` // e.g. "npm", "pypi", "go"
	Version  string `yaml:"version"`
}

// Notifiers holds configuration for supported notification channels.
type Notifiers struct {
	Slack SlackConfig `yaml:"slack"`
	Email EmailConfig `yaml:"email"`
}

// SlackConfig holds Slack webhook settings.
type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`
}

// EmailConfig holds SMTP settings for email notifications.
type EmailConfig struct {
	SMTPHost   string   `yaml:"smtp_host"`
	SMTPPort   int      `yaml:"smtp_port"`
	From       string   `yaml:"from"`
	To         []string `yaml:"to"`
	Password   string   `yaml:"password"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation failed: %w", err)
	}

	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = 24 * time.Hour
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Dependencies) == 0 {
		return fmt.Errorf("at least one dependency must be specified")
	}
	for i, dep := range c.Dependencies {
		if dep.Name == "" {
			return fmt.Errorf("dependency[%d]: name is required", i)
		}
		if dep.Ecosystem == "" {
			return fmt.Errorf("dependency[%d] %q: ecosystem is required", i, dep.Name)
		}
	}
	return nil
}
