package changelog

import "strings"

// Suppressor drops entries whose dependency name matches a suppression list.
// It is useful for silencing known-noisy packages during digest generation.
type Suppressor struct {
	suppressed map[string]struct{}
}

// NewSuppressor returns a Suppressor that will drop entries belonging to any
// dependency in the provided list. Matching is case-insensitive.
func NewSuppressor(deps []string) *Suppressor {
	m := make(map[string]struct{}, len(deps))
	for _, d := range deps {
		m[strings.ToLower(strings.TrimSpace(d))] = struct{}{}
	}
	return &Suppressor{suppressed: m}
}

// Apply removes entries whose Dependency field is in the suppression list.
func (s *Suppressor) Apply(entries []Entry) []Entry {
	if len(s.suppressed) == 0 {
		return entries
	}
	out := entries[:0:0]
	for _, e := range entries {
		key := strings.ToLower(strings.TrimSpace(e.Dependency))
		if _, ok := s.suppressed[key]; !ok {
			out = append(out, e)
		}
	}
	return out
}

// IsSuppressed reports whether a dependency name is currently suppressed.
func (s *Suppressor) IsSuppressed(dep string) bool {
	_, ok := s.suppressed[strings.ToLower(strings.TrimSpace(dep))]
	return ok
}
