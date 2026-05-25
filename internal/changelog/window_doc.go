// Package changelog provides types and utilities for fetching, parsing,
// filtering, and enriching dependency changelog entries.
//
// Window
//
// Window filters a slice of Entry values to those whose Date falls within a
// caller-specified half-open time interval [Start, End). Either bound may be
// left as the zero value to indicate "no lower bound" or "no upper bound"
// respectively.
//
// Example usage:
//
//	w, err := changelog.NewWindow(startTime, endTime)
//	if err != nil {
//		log.Fatal(err)
//	}
//	filtered := w.Apply(entries)
package changelog
