package changelog

import "strings"

// AlertRule defines a condition and severity for triggering an alert on an Entry.
type AlertRule struct {
	Tag      string
	Severity string
}

// Alert holds the result of evaluating an entry against alert rules.
type Alert struct {
	Entry    Entry
	Severity string
	Reason   string
}

// Alerter evaluates entries against a set of rules and produces alerts.
type Alerter struct {
	rules []AlertRule
}

// AlerterOption configures an Alerter.
type AlerterOption func(*Alerter)

// WithAlertRule adds a rule that fires when an entry carries the given tag.
func WithAlertRule(tag, severity string) AlerterOption {
	return func(a *Alerter) {
		a.rules = append(a.rules, AlertRule{Tag: strings.ToLower(tag), Severity: severity})
	}
}

// NewAlerter constructs an Alerter with the provided options.
func NewAlerter(opts ...AlerterOption) *Alerter {
	a := &Alerter{}
	for _, o := range opts {
		o(a)
	}
	return a
}

// Evaluate inspects each entry and returns any alerts triggered by the rules.
func (a *Alerter) Evaluate(entries []Entry) []Alert {
	var alerts []Alert
	for _, e := range entries {
		for _, rule := range a.rules {
			if hasTag(e.Tags, rule.Tag) {
				alerts = append(alerts, Alert{
					Entry:    e,
					Severity: rule.Severity,
					Reason:   "tag:" + rule.Tag,
				})
				break
			}
		}
	}
	return alerts
}

func hasTag(tags []string, target string) bool {
	for _, t := range tags {
		if strings.ToLower(t) == target {
			return true
		}
	}
	return false
}
