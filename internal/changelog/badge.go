package changelog

// Badge represents a visual indicator attached to a changelog entry,
// summarising its importance or category at a glance.
type Badge struct {
	Label string
	Color string // hex colour, e.g. "#e11d48"
}

// BadgeRule maps a tag or label to a Badge.
type BadgeRule struct {
	Tag   string
	Badge Badge
}

// Badger attaches Badge values to entries based on their Tags.
type Badger struct {
	rules []BadgeRule
}

// BadgerOption is a functional option for Badger.
type BadgerOption func(*Badger)

// WithBadgeRule adds a rule that maps tag to the given badge.
func WithBadgeRule(tag, label, color string) BadgerOption {
	return func(b *Badger) {
		b.rules = append(b.rules, BadgeRule{
			Tag:   tag,
			Badge: Badge{Label: label, Color: color},
		})
	}
}

// NewBadger constructs a Badger with the supplied options.
func NewBadger(opts ...BadgerOption) *Badger {
	b := &Badger{}
	for _, o := range opts {
		o(b)
	}
	return b
}

// Apply iterates over entries and attaches badges whose tag matches any
// tag already present on the entry. Existing badges are preserved.
func (b *Badger) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		for _, rule := range b.rules {
			if containsTag(e.Tags, rule.Tag) {
				e.Badges = appendUniqueBadge(e.Badges, rule.Badge)
			}
		}
		out[i] = e
	}
	return out
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func appendUniqueBadge(badges []Badge, b Badge) []Badge {
	for _, existing := range badges {
		if existing.Label == b.Label {
			return badges
		}
	}
	return append(badges, b)
}
