package changelog

// SourceType identifies the kind of changelog source.
type SourceType string

const (
	// SourceHTTP fetches a raw changelog over HTTP/HTTPS.
	SourceHTTP SourceType = "http"
	// SourceGitHub fetches release notes from the GitHub Releases API.
	SourceGitHub SourceType = "github"
)

// Source describes a single dependency changelog source as declared in
// the configuration file.
type Source struct {
	// Name is the human-readable dependency label (e.g. "cobra").
	Name string
	// Type selects the fetching strategy.
	Type SourceType
	// URL is used when Type == SourceHTTP.
	URL string
	// Owner is the GitHub organisation or user; used when Type == SourceGitHub.
	Owner string
	// Repo is the GitHub repository name; used when Type == SourceGitHub.
	Repo string
}

// Validate returns a non-nil error when the source is misconfigured.
func (s Source) Validate() error {
	if s.Name == "" {
		return ErrMissingName
	}
	switch s.Type {
	case SourceHTTP:
		if s.URL == "" {
			return ErrMissingURL
		}
	case SourceGitHub:
		if s.Owner == "" {
			return ErrMissingOwner
		}
		if s.Repo == "" {
			return ErrMissingRepo
		}
	default:
		return ErrUnknownSourceType
	}
	return nil
}

// String returns a short human-readable representation of the source,
// useful for logging and diagnostic output.
func (s Source) String() string {
	switch s.Type {
	case SourceGitHub:
		return string(s.Type) + ":" + s.Owner + "/" + s.Repo
	case SourceHTTP:
		return string(s.Type) + ":" + s.URL
	default:
		return string(s.Type) + ":" + s.Name
	}
}
