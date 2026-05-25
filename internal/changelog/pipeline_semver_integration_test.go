package changelog_test

import (
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
)

func makeSemVerEntries() []changelog.Entry {
	now := time.Now()
	return []changelog.Entry{
		{Dependency: "dep", Version: "v0.9.0", Date: now},
		{Dependency: "dep", Version: "v1.0.0", Date: now},
		{Dependency: "dep", Version: "v1.5.0", Date: now},
		{Dependency: "dep", Version: "v2.0.0", Date: now},
		{Dependency: "dep", Version: "v2.1.0", Date: now},
	}
}

func TestPipeline_SemVerFilter_RemovesOutOfRange(t *testing.T) {
	entries := makeSemVerEntries()
	filter := changelog.NewSemVerFilter(
		changelog.WithMinVersion("v1.0.0"),
		changelog.WithMaxVersion("v2.0.0"),
	)
	out := filter.Apply(entries)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries in range, got %d", len(out))
	}
	for _, e := range out {
		if e.Version == "v0.9.0" || e.Version == "v2.1.0" {
			t.Errorf("out-of-range version %q should have been removed", e.Version)
		}
	}
}

func TestPipeline_SemVerFilter_AllInRange_NoneDropped(t *testing.T) {
	entries := makeSemVerEntries()
	filter := changelog.NewSemVerFilter(
		changelog.WithMinVersion("v0.1.0"),
		changelog.WithMaxVersion("v9.0.0"),
	)
	out := filter.Apply(entries)
	if len(out) != len(entries) {
		t.Fatalf("expected all %d entries, got %d", len(entries), len(out))
	}
}

func TestPipeline_SemVerFilter_EmptyInput_DoesNotPanic(t *testing.T) {
	filter := changelog.NewSemVerFilter(changelog.WithMinVersion("v1.0.0"))
	out := filter.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(out))
	}
}
