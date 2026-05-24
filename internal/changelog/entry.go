package changelog

import "time"

// Entry represents a single parsed changelog or release entry.
type Entry struct {
	Dependency string
	Version    string
	Date       time.Time
	Body       string
	Link       string
	Tags       []string
	Badges     []Badge
	Labels     []string
	Score      int
	Highlighted bool
	Summary    string
	Category   string
	Refs       []string
}
