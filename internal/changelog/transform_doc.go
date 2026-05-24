// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelogs.
//
// The Transformer interface and Chain type allow composable, ordered
// transformations over []Entry slices. Use NewChain to combine multiple
// Transformer steps into a single pass, and LimitTransformer to cap
// the number of entries returned to downstream consumers.
package changelog
