// Package store implements a lightweight file-backed key-value store used by
// depwatch to remember the last changelog version seen for each monitored
// dependency.
//
// The store serialises its state as a JSON object on disk so that version
// history survives daemon restarts.  All public methods are safe for
// concurrent use.
package store
