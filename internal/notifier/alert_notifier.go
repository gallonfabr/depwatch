package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourorg/depwatch/internal/changelog"
)

// AlertNotifier wraps an existing Notifier and sends a message
// whenever the provided Alerter fires on a set of entries.
type AlertNotifier struct {
	inner   Notifier
	alerter *changelog.Alerter
}

// NewAlertNotifier constructs an AlertNotifier.
// Returns an error if inner or alerter is nil.
func NewAlertNotifier(inner Notifier, alerter *changelog.Alerter) (*AlertNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("inner notifier must not be nil")
	}
	if alerter == nil {
		return nil, fmt.Errorf("alerter must not be nil")
	}
	return &AlertNotifier{inner: inner, alerter: alerter}, nil
}

// Notify evaluates entries and sends one consolidated message per severity
// group if any alerts are triggered. Returns nil when no alerts fire.
func (an *AlertNotifier) Notify(ctx context.Context, entries []changelog.Entry) error {
	alerts := an.alerter.Evaluate(entries)
	if len(alerts) == 0 {
		return nil
	}

	bySeverity := make(map[string][]string)
	for _, a := range alerts {
		bySeverity[a.Severity] = append(bySeverity[a.Severity],
			fmt.Sprintf("%s@%s (%s)", a.Entry.Dependency, a.Entry.Version, a.Reason))
	}

	var sb strings.Builder
	sb.WriteString("[depwatch] Alerts fired:\n")
	for sev, items := range bySeverity {
		sb.WriteString(fmt.Sprintf("  [%s]\n", strings.ToUpper(sev)))
		for _, item := range items {
			sb.WriteString("    - " + item + "\n")
		}
	}

	return an.inner.Send(ctx, sb.String())
}
