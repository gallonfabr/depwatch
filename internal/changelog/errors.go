package changelog

import "errors"

// Sentinel errors shared across the changelog package.
var (
	// ErrMissingURL is returned when an HTTP source has no URL configured.
	ErrMissingURL = errors.New("changelog: http source requires a url")

	// ErrMissingOwner is returned when a GitHub source has no owner configured.
	ErrMissingOwner = errors.New("changelog: github source requires an owner")

	// ErrMissingRepo is returned when a GitHub source has no repo configured.
	ErrMissingRepo = errors.New("changelog: github source requires a repo")

	// ErrUnknownSourceType is returned when the source type is not recognised.
	ErrUnknownSourceType = errors.New("changelog: unknown source type")
)
