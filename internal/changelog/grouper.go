package changelog

// Group holds a label and the entries assigned to that label.
type Group struct {
	Label   string
	Entries []Entry
}

// Grouper partitions a slice of entries by their first label.
// Entries with no labels are placed under the fallback label.
type Grouper struct {
	fallback string
	order    []string
}

// GrouperOption configures a Grouper.
type GrouperOption func(*Grouper)

// WithFallbackLabel sets the label used for entries that carry no labels.
func WithFallbackLabel(label string) GrouperOption {
	return func(g *Grouper) {
		if label != "" {
			g.fallback = label
		}
	}
}

// WithGroupOrder defines the preferred display order of label groups.
// Groups not listed here appear after the ordered ones.
func WithGroupOrder(labels ...string) GrouperOption {
	return func(g *Grouper) {
		g.order = labels
	}
}

// NewGrouper returns a Grouper with optional configuration applied.
func NewGrouper(opts ...GrouperOption) *Grouper {
	g := &Grouper{fallback: "other"}
	for _, o := range opts {
		o(g)
	}
	return g
}

// Apply groups entries by their first label and returns ordered groups.
func (g *Grouper) Apply(entries []Entry) []Group {
	buckets := make(map[string][]Entry)
	for _, e := range entries {
		key := g.fallback
		if len(e.Labels) > 0 {
			key = e.Labels[0]
		}
		buckets[key] = append(buckets[key], e)
	}

	seen := make(map[string]bool)
	var groups []Group

	for _, label := range g.order {
		if entries, ok := buckets[label]; ok {
			groups = append(groups, Group{Label: label, Entries: entries})
			seen[label] = true
		}
	}

	for key, entries := range buckets {
		if !seen[key] {
			groups = append(groups, Group{Label: key, Entries: entries})
		}
	}

	return groups
}
