package changelog

// Annotator attaches arbitrary key-value metadata to changelog entries.
// Annotations are stored in the entry's Tags slice using a "key:value" format
// so they remain compatible with the rest of the pipeline.
type Annotator struct {
	annotations map[string]string
}

// AnnotatorOption configures an Annotator.
type AnnotatorOption func(*Annotator)

// WithAnnotation registers a single key-value annotation that will be
// appended to every entry processed by the Annotator.
func WithAnnotation(key, value string) AnnotatorOption {
	return func(a *Annotator) {
		if key != "" {
			a.annotations[key] = value
		}
	}
}

// NewAnnotator constructs an Annotator with the supplied options.
func NewAnnotator(opts ...AnnotatorOption) *Annotator {
	a := &Annotator{
		annotations: make(map[string]string),
	}
	for _, o := range opts {
		o(a)
	}
	return a
}

// Apply attaches all registered annotations to each entry and returns the
// modified slice. The original slice is not mutated; each entry is copied.
func (a *Annotator) Apply(entries []Entry) []Entry {
	if len(a.annotations) == 0 {
		return entries
	}
	out := make([]Entry, len(entries))
	for i, e := range entries {
		copy_ := e
		for k, v := range a.annotations {
			tag := k + ":" + v
			if !annotatorHasTag(copy_.Tags, tag) {
				copy_.Tags = append(copy_.Tags, tag)
			}
		}
		out[i] = copy_
	}
	return out
}

func annotatorHasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
