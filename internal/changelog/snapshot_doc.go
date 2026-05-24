// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelogs.
//
// SnapshotStore offers a concurrency-safe, in-memory store for the most
// recent set of changelog entries per dependency. It is intended to be used
// by the watcher layer to detect changes between polling cycles without
// requiring a full diff against persistent storage on every tick.
//
// Typical usage:
//
//	store := changelog.NewSnapshotStore()
//	store.Save("github.com/some/dep", entries)
//	if snap, ok := store.Get("github.com/some/dep"); ok {
//		// compare snap.Entries with newly fetched entries
//	}
package changelog
