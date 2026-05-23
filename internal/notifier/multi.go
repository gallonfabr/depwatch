// Package notifier provides implementations for sending digest notifications.
package notifier

import (
	"errors"
	"fmt"
)

// Notifier defines the interface for sending a digest message.
type Notifier interface {
	Send(subject, body string) error
}

// MultiNotifier fans out a single Send call to multiple Notifier implementations.
// All notifiers are attempted; errors are collected and joined.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a MultiNotifier from the provided notifiers.
// It returns an error if the slice is empty.
func NewMultiNotifier(notifiers ...Notifier) (*MultiNotifier, error) {
	if len(notifiers) == 0 {
		return nil, errors.New("multinotifier: at least one notifier is required")
	}
	return &MultiNotifier{notifiers: notifiers}, nil
}

// Send delivers subject and body to every registered notifier.
// If one or more notifiers fail, their errors are combined and returned.
func (m *MultiNotifier) Send(subject, body string) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Send(subject, body); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("multinotifier: %d notifier(s) failed: %w", len(errs), errors.Join(errs...))
}
