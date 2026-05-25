// Package changelog provides the annotator component which enriches
// changelog entries with arbitrary key-value metadata.
//
// # Annotator
//
// An Annotator attaches static key-value pairs to every entry that passes
// through the pipeline. Annotations are encoded as "key:value" strings and
// stored in the entry's Tags field, keeping them compatible with downstream
// pipeline stages such as the Tagger, Labeler, and Router.
//
// Typical use-cases include stamping entries with deployment environment,
// owning team, or data-source identifiers so that routing and alerting
// rules can filter on those values.
//
//	ann := changelog.NewAnnotator(
//	    changelog.WithAnnotation("env", "production"),
//	    changelog.WithAnnotation("team", "platform"),
//	)
//	enriched := ann.Apply(entries)
package changelog
