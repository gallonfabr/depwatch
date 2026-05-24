package changelog

// Stage is a processing step that transforms a slice of Entry values.
type Stage interface {
	Apply(entries []Entry) []Entry
}

// Pipeline runs a sequence of Stages over a slice of changelog entries.
type Pipeline struct {
	stages []Stage
}

// NewPipeline constructs a Pipeline with the provided stages executed in order.
func NewPipeline(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Run passes entries through each stage in sequence and returns the result.
func (p *Pipeline) Run(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	for _, s := range p.stages {
		out = s.Apply(out)
	}
	return out
}
