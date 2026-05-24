package changelog

// Classifier assigns a category string to each Entry based on its labels.
// It maps known labels to canonical category names and falls back to a
// configurable default when no label matches.
type Classifier struct {
	mappings map[string]string
	fallback string
}

// ClassifierOption configures a Classifier.
type ClassifierOption func(*Classifier)

// WithCategoryMapping registers a mapping from a label to a category name.
func WithCategoryMapping(label, category string) ClassifierOption {
	return func(c *Classifier) {
		c.mappings[label] = category
	}
}

// WithClassifierFallback sets the default category used when no label matches.
func WithClassifierFallback(category string) ClassifierOption {
	return func(c *Classifier) {
		c.fallback = category
	}
}

// NewClassifier returns a Classifier with the provided options applied.
// Default mappings cover common changelog labels. The default fallback
// category is "other".
func NewClassifier(opts ...ClassifierOption) *Classifier {
	c := &Classifier{
		mappings: map[string]string{
			"security": "Security",
			"feature":  "Features",
			"bugfix":   "Bug Fixes",
			"breaking": "Breaking Changes",
			"docs":     "Documentation",
		},
		fallback: "Other",
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Apply iterates over entries and sets Entry.Category based on the first
// matching label found in the classifier's mappings. If no label matches,
// the fallback category is used.
func (c *Classifier) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		category := c.fallback
		for _, lbl := range e.Labels {
			if cat, ok := c.mappings[lbl]; ok {
				category = cat
				break
			}
		}
		e.Category = category
		out[i] = e
	}
	return out
}
