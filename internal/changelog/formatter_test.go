package changelog

import (
	"strings"
	"testing"
	"time"
)

func sampleFormatterEntries() []Entry {
	return []Entry{
		{
			Dependency: "github.com/foo/bar",
			Version:    "v1.2.0",
			Date:       time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			Labels:     []string{"feature"},
			Badges:     []string{"🚀"},
			Body:       "Added new API.",
		},
		{
			Dependency: "github.com/foo/bar",
			Version:    "v1.1.0",
			Date:       time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			Labels:     []string{"bugfix"},
			Body:       "",
		},
		{
			Dependency: "github.com/baz/qux",
			Version:    "v0.9.0",
			Date:       time.Time{},
			Labels:     []string{},
		},
	}
}

func TestNewFormatter_NotNil(t *testing.T) {
	f := NewFormatter()
	if f == nil {
		t.Fatal("expected non-nil Formatter")
	}
}

func TestFormatter_Format_Empty(t *testing.T) {
	f := NewFormatter()
	if got := f.Format(nil); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestFormatter_Format_ContainsDependencyHeader(t *testing.T) {
	f := NewFormatter()
	out := f.Format(sampleFormatterEntries())
	if !strings.Contains(out, "## github.com/foo/bar") {
		t.Errorf("expected dependency header in output")
	}
	if !strings.Contains(out, "## github.com/baz/qux") {
		t.Errorf("expected second dependency header in output")
	}
}

func TestFormatter_Format_ContainsVersion(t *testing.T) {
	f := NewFormatter()
	out := f.Format(sampleFormatterEntries())
	if !strings.Contains(out, "v1.2.0") {
		t.Errorf("expected version in output")
	}
}

func TestFormatter_Format_ContainsDate(t *testing.T) {
	f := NewFormatter()
	out := f.Format(sampleFormatterEntries())
	if !strings.Contains(out, "2024-03-15") {
		t.Errorf("expected formatted date in output")
	}
}

func TestFormatter_Format_ZeroDate_Omitted(t *testing.T) {
	f := NewFormatter()
	out := f.Format(sampleFormatterEntries())
	// entry for baz/qux has zero date; line should not contain a date token
	lines := strings.Split(out, "\n")
	for _, l := range lines {
		if strings.Contains(l, "v0.9.0") && strings.Contains(l, "0001-") {
			t.Errorf("zero date should not appear in output, got: %s", l)
		}
	}
}

func TestFormatter_Format_ShowLabels(t *testing.T) {
	f := NewFormatter(WithLabelsEnabled(true))
	out := f.Format(sampleFormatterEntries())
	if !strings.Contains(out, "[feature]") {
		t.Errorf("expected label in output")
	}
}

func TestFormatter_Format_ShowBadges(t *testing.T) {
	f := NewFormatter(WithBadgesEnabled(true))
	out := f.Format(sampleFormatterEntries())
	if !strings.Contains(out, "🚀") {
		t.Errorf("expected badge in output")
	}
}

func TestFormatter_Format_CustomDateLayout(t *testing.T) {
	f := NewFormatter(WithDateLayout("01/02/2006"))
	out := f.Format(sampleFormatterEntries())
	if !strings.Contains(out, "03/15/2024") {
		t.Errorf("expected custom date layout in output, got:\n%s", out)
	}
}
