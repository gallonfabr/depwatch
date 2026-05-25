package changelog

import (
	"fmt"
	"strings"
	"time"
)

// Formatter converts a slice of Entry values into a human-readable string
// suitable for inclusion in a digest notification.
type Formatter struct {
	dateLayout string
	showBadges bool
	showLabels bool
}

// FormatterOption configures a Formatter.
type FormatterOption func(*Formatter)

// WithDateLayout overrides the date layout used when rendering entry dates.
func WithDateLayout(layout string) FormatterOption {
	return func(f *Formatter) {
		if layout != "" {
			f.dateLayout = layout
		}
	}
}

// WithBadgesEnabled controls whether badges are appended to each entry line.
func WithBadgesEnabled(enabled bool) FormatterOption {
	return func(f *Formatter) { f.showBadges = enabled }
}

// WithLabelsEnabled controls whether labels are appended to each entry line.
func WithLabelsEnabled(enabled bool) FormatterOption {
	return func(f *Formatter) { f.showLabels = enabled }
}

// NewFormatter creates a Formatter with sensible defaults.
func NewFormatter(opts ...FormatterOption) *Formatter {
	f := &Formatter{
		dateLayout: "2006-01-02",
		showBadges: false,
		showLabels: true,
	}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Format renders entries into a multi-line string grouped by dependency.
func (f *Formatter) Format(entries []Entry) string {
	if len(entries) == 0 {
		return ""
	}

	var sb strings.Builder
	current := ""

	for _, e := range entries {
		if e.Dependency != current {
			if current != "" {
				sb.WriteByte('\n')
			}
			fmt.Fprintf(&sb, "## %s\n", e.Dependency)
			current = e.Dependency
		}

		date := ""
		if !e.Date.IsZero() {
			date = e.Date.UTC().Format(f.dateLayout) + " "
		}

		line := fmt.Sprintf("- %s%s", date, e.Version)

		if f.showLabels && len(e.Labels) > 0 {
			line += fmt.Sprintf(" [%s]", strings.Join(e.Labels, ", "))
		}
		if f.showBadges && len(e.Badges) > 0 {
			line += fmt.Sprintf(" (%s)", strings.Join(e.Badges, " "))
		}
		if e.Body != "" {
			line += "\n  " + strings.ReplaceAll(strings.TrimSpace(e.Body), "\n", "\n  ")
		}

		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return sb.String()
}

// ensure time is imported only when needed by tests via the zero-value check
var _ = time.Time{}
