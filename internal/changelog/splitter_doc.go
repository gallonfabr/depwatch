// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelog entries.
//
// Splitter
//
// Splitter partitions a flat []Entry slice into per-dependency buckets,
// making it straightforward to process or render each dependency's
// changelog independently.
//
// Example:
//
//	s := changelog.NewSplitter()
//	buckets := s.Split(entries)
//	for dep, depEntries := range buckets {
//		fmt.Println(dep, len(depEntries))
//	}
package changelog
