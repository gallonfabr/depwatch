package changelog_test

import (
	"testing"

	"github.com/yourorg/depwatch/internal/changelog"
)

func TestSource_Validate_HTTPMissingURL(t *testing.T) {
	s := changelog.Source{Name: "lib", Type: changelog.SourceHTTP}
	if err := s.Validate(); err != changelog.ErrMissingURL {
		t.Fatalf("expected ErrMissingURL, got %v", err)
	}
}

func TestSource_Validate_HTTPValid(t *testing.T) {
	s := changelog.Source{Name: "lib", Type: changelog.SourceHTTP, URL: "https://example.com/CHANGELOG.md"}
	if err := s.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSource_Validate_GitHubMissingOwner(t *testing.T) {
	s := changelog.Source{Name: "lib", Type: changelog.SourceGitHub, Repo: "myrepo"}
	if err := s.Validate(); err != changelog.ErrMissingOwner {
		t.Fatalf("expected ErrMissingOwner, got %v", err)
	}
}

func TestSource_Validate_GitHubMissingRepo(t *testing.T) {
	s := changelog.Source{Name: "lib", Type: changelog.SourceGitHub, Owner: "myorg"}
	if err := s.Validate(); err != changelog.ErrMissingRepo {
		t.Fatalf("expected ErrMissingRepo, got %v", err)
	}
}

func TestSource_Validate_GitHubValid(t *testing.T) {
	s := changelog.Source{Name: "lib", Type: changelog.SourceGitHub, Owner: "myorg", Repo: "myrepo"}
	if err := s.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSource_Validate_UnknownType(t *testing.T) {
	s := changelog.Source{Name: "lib", Type: changelog.SourceType("rss")}
	if err := s.Validate(); err != changelog.ErrUnknownSourceType {
		t.Fatalf("expected ErrUnknownSourceType, got %v", err)
	}
}

func TestSourceType_Constants(t *testing.T) {
	if changelog.SourceHTTP != "http" {
		t.Errorf("SourceHTTP should be 'http'")
	}
	if changelog.SourceGitHub != "github" {
		t.Errorf("SourceGitHub should be 'github'")
	}
}
