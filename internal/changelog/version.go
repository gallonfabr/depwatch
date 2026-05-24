package changelog

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a parsed semantic version.
type Version struct {
	Major int
	Minor int
	Patch int
	Pre   string
}

// ParseVersion parses a version string like "v1.2.3" or "1.2.3-beta".
// It returns an error if the version string is malformed.
func ParseVersion(s string) (Version, error) {
	s = strings.TrimPrefix(s, "v")

	pre := ""
	if idx := strings.IndexByte(s, '-'); idx != -1 {
		pre = s[idx+1:]
		s = s[:idx]
	}

	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("changelog: invalid version %q: expected major.minor.patch", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, fmt.Errorf("changelog: invalid major in %q: %w", s, err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("changelog: invalid minor in %q: %w", s, err)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, fmt.Errorf("changelog: invalid patch in %q: %w", s, err)
	}

	return Version{Major: major, Minor: minor, Patch: patch, Pre: pre}, nil
}

// String returns the canonical string representation of the version.
func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	return s
}

// Less reports whether v is older than other.
func (v Version) Less(other Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	return v.Patch < other.Patch
}

// Equal reports whether v and other represent the same version,
// ignoring pre-release labels.
func (v Version) Equal(other Version) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}
