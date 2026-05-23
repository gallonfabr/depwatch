// Package main is the entry point for the depwatch daemon.
//
// depwatch monitors dependency changelogs at a configured interval and
// delivers digest notifications via Slack or email.
//
// Usage:
//
//	depwatch [config-path]
//
// If no config path is provided, depwatch.yaml in the current directory
// is used by default.
//
// The daemon handles SIGINT and SIGTERM for graceful shutdown.
package main
