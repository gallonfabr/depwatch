// Package notifier provides implementations for delivering dependency digest
// notifications to external services.
//
// Supported notifiers:
//
//   - SlackNotifier: posts messages to a Slack incoming webhook URL.
//   - EmailNotifier: sends messages via SMTP using optional PlainAuth.
//
// Each notifier is constructed via a New* constructor that validates required
// configuration at creation time, returning a descriptive error when fields
// are missing or invalid.
package notifier
