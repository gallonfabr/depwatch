package changelog_test

import (
	"testing"
	"time"

	"github.com/depwatch/internal/changelog"
)

var pipelineEntries = []changelog.Entry{
	{Dependency: "react", Version: "18.2.0", Date: time.Now().Add(-1 * time.Hour), Body: "<b>Fix</b> hydration bug"},
	{Dependency: "react", Version: "18.1.0", Date: time.Now().Add(-48 * time.Hour), Body: "Minor update"},
	{Dependency: "react", Version: "18.2.0", Date: time.Now().Add(-1 * time.Hour), Body: "<b>Fix</b> hydration bug"},
	{Dependency: "lodash", Version: "4.17.21", Date: time.Now().Add(-2 * time.Hour), Body: "Security patch"},
}

func TestNewPipeline_NotNil(t *testing.T) {
	p := changelog.NewPipeline()
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestPipeline_Run_Empty(t *testing.T) {
	p := changelog.NewPipeline()
	result := p.Run([]changelog.Entry{})
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestPipeline_Run_DeduplicatesEntries(t *testing.T) {
	p := changelog.NewPipeline(
		changelog.NewDeduplicator(),
	)
	result := p.Run(pipelineEntries)
	// react 18.2.0 appears twice; should be deduplicated
	count := 0
	for _, e := range result {
		if e.Dependency == "react" && e.Version == "18.2.0" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 deduplicated react 18.2.0 entry, got %d", count)
	}
}

func TestPipeline_Run_NormalizesAndDedupes(t *testing.T) {
	p := changelog.NewPipeline(
		changelog.NewDeduplicator(),
		changelog.NewNormalizer(changelog.WithStripHTML()),
	)
	result := p.Run(pipelineEntries)
	for _, e := range result {
		if e.Body == "<b>Fix</b> hydration bug" {
			t.Errorf("expected HTML to be stripped, got: %s", e.Body)
		}
	}
}

func TestPipeline_Run_SortsDescending(t *testing.T) {
	p := changelog.NewPipeline(
		changelog.NewSorter(false),
	)
	result := p.Run(pipelineEntries)
	for i := 1; i < len(result); i++ {
		if result[i].Date.After(result[i-1].Date) {
			t.Errorf("entries not sorted descending at index %d", i)
		}
	}
}

func TestPipeline_Run_FullChain(t *testing.T) {
	p := changelog.NewPipeline(
		changelog.NewDeduplicator(),
		changelog.NewNormalizer(changelog.WithStripHTML()),
		changelog.NewSorter(false),
	)
	result := p.Run(pipelineEntries)
	// 4 entries minus 1 duplicate = 3
	if len(result) != 3 {
		t.Errorf("expected 3 entries after full pipeline, got %d", len(result))
	}
}
