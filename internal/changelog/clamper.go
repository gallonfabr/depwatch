package changelog

// Clamper ensures that numeric score values on entries stay within a
// defined [Min, Max] range. Entries whose scores already fall within
// the range are left unchanged.
type Clamper struct {
	min float64
	max float64
}

// ClamperOption configures a Clamper.
type ClamperOption func(*Clamper)

// WithClampMin sets the minimum allowed score.
func WithClampMin(min float64) ClamperOption {
	return func(c *Clamper) {
		c.min = min
	}
}

// WithClampMax sets the maximum allowed score.
func WithClampMax(max float64) ClamperOption {
	return func(c *Clamper) {
		c.max = max
	}
}

// NewClamper returns a Clamper with the supplied options applied.
// If no options are provided the default range is [0, 100].
func NewClamper(opts ...ClamperOption) *Clamper {
	c := &Clamper{
		min: 0,
		max: 100,
	}
	for _, o := range opts {
		o(c)
	}
	if c.min > c.max {
		c.min, c.max = c.max, c.min
	}
	return c
}

// Apply clamps the Score field of every entry to the configured range
// and returns the modified slice. The input slice is mutated in place.
func (c *Clamper) Apply(entries []Entry) []Entry {
	for i := range entries {
		if entries[i].Score < c.min {
			entries[i].Score = c.min
		}
		if entries[i].Score > c.max {
			entries[i].Score = c.max
		}
	}
	return entries
}
