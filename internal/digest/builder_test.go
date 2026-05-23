package digest

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
)

func sampleUpdates() map[string][]changelog.Entry {
	return map[string][]changelog.Entry{
		"cobra": {
			{
				Version: "v1.8.0",
				Date:    time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
				Body:    "Added shell completion improvements.",
			},
		},
		"viper": {
			{
				Version: "v1.18.0",
				Date:    time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
				Body:    "Fixed env var binding edge case.",
			},
		},
	}
}

func TestBuild_EntryCount(t *testing.T) {
	b := NewBuilder()
	d := b.Build(sampleUpdates())
	if len(d.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(d.Entries))
	}
}

func TestBuild_EmptyUpdates(t *testing.T) {
	b := NewBuilder()
	d := b.Build(map[string][]changelog.Entry{})
	if len(d.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(d.Entries))
	}
}

func TestFormatText_ContainsDependencyName(t *testing.T) {
	b := NewBuilder()
	d := b.Build(sampleUpdates())
	text := d.FormatText()
	if !strings.Contains(text, "cobra") && !strings.Contains(text, "viper") {
		t.Error("expected formatted text to contain dependency names")
	}
}

func TestFormatText_EmptyDigest(t *testing.T) {
	d := &Digest{}
	text := d.FormatText()
	if text != "No dependency updates found." {
		t.Errorf("unexpected empty digest text: %q", text)
	}
}

func TestFormatText_ContainsVersion(t *testing.T) {
	b := NewBuilder()
	d := b.Build(map[string][]changelog.Entry{
		"cobra": {{Version: "v1.8.0", Date: time.Now(), Body: "some notes"}},
	})
	text := d.FormatText()
	if !strings.Contains(text, "v1.8.0") {
		t.Errorf("expected formatted text to contain version, got:\n%s", text)
	}
}

func TestBuild_GeneratedAtSet(t *testing.T) {
	b := NewBuilder()
	before := time.Now().UTC()
	d := b.Build(sampleUpdates())
	after := time.Now().UTC()
	if d.GeneratedAt.Before(before) || d.GeneratedAt.After(after) {
		t.Error("GeneratedAt timestamp is out of expected range")
	}
}
