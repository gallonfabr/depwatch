// Package changelog provides primitives for fetching, parsing, and processing
// dependency changelog entries.
//
// # Classifier
//
// Classifier assigns a human-readable category to each Entry by inspecting its
// Labels field. Built-in mappings cover the most common changelog categories:
//
//	"security"  → "Security"
//	"feature"   → "Features"
//	"bugfix"    → "Bug Fixes"
//	"breaking"  → "Breaking Changes"
//	"docs"      → "Documentation"
//
// Custom mappings can be added via WithCategoryMapping and the fallback
// category (default "Other") can be overridden with WithClassifierFallback.
package changelog
