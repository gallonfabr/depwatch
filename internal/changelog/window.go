package changelog

import (
	"errors"
	"time"
)

// ErrInvalidWindow is returned when a Window is misconfigured.
var ErrInvalidWindow = errors.New("changelog: window start must be before end")

// Window restricts a slice of Entry values to those whose Date falls within
// a half-open interval [Start, End). A zero-value Start or End disables that
// bound.
type Window struct {
	Start time.Time
	End   time.Time
}

// NewWindow returns a Window for the given bounds. It returns
// ErrInvalidWindow when both bounds are non-zero and Start is not before End.
func NewWindow(start, end time.Time) (Window, error) {
	if !start.IsZero() && !end.IsZero() && !start.Before(end) {
		return Window{}, ErrInvalidWindow
	}
	return Window{Start: start, End: end}, nil
}

// Apply returns the subset of entries whose Date falls within the window.
func (w Window) Apply(entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !w.Start.IsZero() && e.Date.Before(w.Start) {
			continue
		}
		if !w.End.IsZero() && !e.Date.Before(w.End) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// IsZero reports whether both bounds are unset.
func (w Window) IsZero() bool {
	return w.Start.IsZero() && w.End.IsZero()
}
