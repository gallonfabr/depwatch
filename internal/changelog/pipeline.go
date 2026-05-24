package changelog

// Pipeline chains multiple entry transformers together, applying
// each stage in order to a slice of changelog entries.
//
// Supported stage types:
//   - *Filter       – filters by date/limit
//   - *Deduplicator – removes duplicate versions
//   - *Normalizer   – normalises body text
//   - *Sorter       – sorts by date
type Pipeline struct {
	stages []stage
}

type stage interface {
	Apply([]Entry) []Entry
}

// NewPipeline constructs a Pipeline from the provided stages.
// Stages are applied in the order they are given.
func NewPipeline(stages ...stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Run executes every stage in sequence and returns the final entries.
func (p *Pipeline) Run(entries []Entry) []Entry {
	out := entries
	for _, s := range p.stages {
		out = s.Apply(out)
	}
	return out
}
