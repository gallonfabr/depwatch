package changelog

// Transformer applies a transformation to a slice of Entry values.
type Transformer interface {
	Transform(entries []Entry) []Entry
}

// TransformFunc is a function adapter for Transformer.
type TransformFunc func(entries []Entry) []Entry

// Transform implements Transformer.
func (f TransformFunc) Transform(entries []Entry) []Entry {
	return f(entries)
}

// Chain combines multiple Transformers into a single Transformer that
// applies each in order.
type Chain struct {
	steps []Transformer
}

// NewChain creates a Chain from the provided Transformers.
// An empty chain is valid and returns entries unchanged.
func NewChain(steps ...Transformer) *Chain {
	return &Chain{steps: steps}
}

// Transform applies each step in sequence to the entries.
func (c *Chain) Transform(entries []Entry) []Entry {
	for _, s := range c.steps {
		entries = s.Transform(entries)
	}
	return entries
}

// LimitTransformer caps the number of returned entries.
type LimitTransformer struct {
	max int
}

// NewLimitTransformer returns a Transformer that keeps at most max entries.
func NewLimitTransformer(max int) *LimitTransformer {
	if max < 0 {
		max = 0
	}
	return &LimitTransformer{max: max}
}

// Transform truncates entries to the configured maximum.
func (l *LimitTransformer) Transform(entries []Entry) []Entry {
	if l.max == 0 || len(entries) <= l.max {
		return entries
	}
	return entries[:l.max]
}
