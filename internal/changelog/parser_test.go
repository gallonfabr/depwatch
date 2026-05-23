package changelog

import (
	"testing"
	"time"
)

const sampleChangelog = `# Changelog

## [1.2.0] - 2024-03-10
### Added
- New feature A
- New feature B

## 1.1.0 (2024-01-05)
### Fixed
- Bug fix X

## v1.0.0 — 2023-11-20
Initial release.
`

func TestParser_Parse_EntryCount(t *testing.T) {
	p := NewParser()
	entries := p.Parse(sampleChangelog)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestParser_Parse_Versions(t *testing.T) {
	p := NewParser()
	entries := p.Parse(sampleChangelog)
	expected := []string{"1.2.0", "1.1.0", "1.0.0"}
	for i, want := range expected {
		if entries[i].Version != want {
			t.Errorf("entry %d: expected version %q, got %q", i, want, entries[i].Version)
		}
	}
}

func TestParser_Parse_Dates(t *testing.T) {
	p := NewParser()
	entries := p.Parse(sampleChangelog)

	wantDate := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
	if !entries[0].Date.Equal(wantDate) {
		t.Errorf("entry 0: expected date %v, got %v", wantDate, entries[0].Date)
	}
	if entries[0].RawDate != "2024-03-10" {
		t.Errorf("entry 0: unexpected RawDate %q", entries[0].RawDate)
	}
}

func TestParser_Parse_Body(t *testing.T) {
	p := NewParser()
	entries := p.Parse(sampleChangelog)
	if entries[2].Body != "Initial release." {
		t.Errorf("entry 2: unexpected body %q", entries[2].Body)
	}
}

func TestParser_Parse_Empty(t *testing.T) {
	p := NewParser()
	entries := p.Parse("")
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries for empty input, got %d", len(entries))
	}
}

func TestParser_Parse_NoVersionHeadings(t *testing.T) {
	p := NewParser()
	entries := p.Parse("Just some random text\nwith no version headings.")
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}
