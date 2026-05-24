package changelog

// Label constants assigned by the Labeler.
const (
	LabelSecurity  = "security"
	LabelBreaking  = "breaking"
	LabelFeature   = "feature"
	LabelBugfix    = "bugfix"
	LabelUnknown   = "unknown"
)

// labelerKeywords maps a label to the keywords that trigger it.
var labelerKeywords = map[string][]string{
	LabelSecurity: {"security", "vulnerability", "cve", "exploit", "patch"},
	LabelBreaking: {"breaking", "incompatible", "removed", "deprecated"},
	LabelFeature:  {"feature", "added", "new", "introduce"},
	LabelBugfix:   {"fix", "bugfix", "bug", "resolved", "patch"},
}

// Labeler assigns a human-readable label to each Entry based on its body text.
type Labeler struct{}

// NewLabeler constructs a Labeler.
func NewLabeler() *Labeler {
	return &Labeler{}
}

// Apply iterates over entries and sets the Label field on each one.
// Labels are evaluated in priority order: security > breaking > feature > bugfix.
// If no keyword matches, LabelUnknown is used.
func (l *Labeler) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		e.Label = l.classify(e.Body)
		out[i] = e
	}
	return out
}

// classify returns the highest-priority label that matches the body text.
func (l *Labeler) classify(body string) string {
	lower := toLower(body)
	priority := []string{LabelSecurity, LabelBreaking, LabelFeature, LabelBugfix}
	for _, label := range priority {
		for _, kw := range labelerKeywords[label] {
			if containsWord(lower, kw) {
				return label
			}
		}
	}
	return LabelUnknown
}

// toLower is a thin wrapper kept here so the file is self-contained.
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// containsWord reports whether substr appears anywhere in s.
func containsWord(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(contains(s, substr)))
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
