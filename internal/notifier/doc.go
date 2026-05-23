// Package notifier provides integrations for delivering digest messages
// to external services such as Slack and email.
//
// Each notifier implements a simple Send(message string) error interface,
// allowing the depwatch daemon to dispatch formatted digests produced by
// the digest package without coupling to any specific delivery mechanism.
//
// Supported notifiers:
//   - SlackNotifier: posts messages to a Slack incoming webhook URL.
package notifier
