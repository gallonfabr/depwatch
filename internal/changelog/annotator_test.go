package changelog

import (
	"strings"
	"testing"
)

func TestNewAnnotator_NotNil(t *testing.T) {
	a := NewAnnotator()
	if a == nil {
		t.Fatal("expected non-nil Annotator")
	}
}

func TestAnnotator_NoOptions_ReturnsSameEntries(t *testing.T) {
	a := NewAnnotator()
	entries := []Entry{{Dependency: "dep-a"}, {Dependency: "dep-b"}}
	out := a.Apply(entries)
	if len(out) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(out))
	}
}

func TestAnnotator_SingleAnnotation_AppliedToAll(t *testing.T) {
	a := NewAnnotator(WithAnnotation("env", "production"))
	entries := []Entry{{Dependency: "dep-a"}, {Dependency: "dep-b"}}
	out := a.Apply(entries)
	for _, e := range out {
		if !annotatorHasTag(e.Tags, "env:production") {
			t.Errorf("entry %q missing annotation tag", e.Dependency)
		}
	}
}

func TestAnnotator_MultipleAnnotations_AllApplied(t *testing.T) {
	a := NewAnnotator(
		WithAnnotation("env", "staging"),
		WithAnnotation("team", "platform"),
	)
	entries := []Entry{{Dependency: "lib"}}
	out := a.Apply(entries)
	if !annotatorHasTag(out[0].Tags, "env:staging") {
		t.Error("missing env:staging tag")
	}
	if !annotatorHasTag(out[0].Tags, "team:platform") {
		t.Error("missing team:platform tag")
	}
}

func TestAnnotator_DoesNotDuplicateTag(t *testing.T) {
	a := NewAnnotator(WithAnnotation("env", "prod"))
	entries := []Entry{{Dependency: "dep", Tags: []string{"env:prod"}}}
	out := a.Apply(entries)
	count := 0
	for _, tag := range out[0].Tags {
		if strings.HasPrefix(tag, "env:") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 env tag, got %d", count)
	}
}

func TestAnnotator_EmptyKey_Ignored(t *testing.T) {
	a := NewAnnotator(WithAnnotation("", "value"))
	entries := []Entry{{Dependency: "dep"}}
	out := a.Apply(entries)
	if len(out[0].Tags) != 0 {
		t.Errorf("expected no tags for empty key, got %v", out[0].Tags)
	}
}

func TestAnnotator_DoesNotMutateInput(t *testing.T) {
	a := NewAnnotator(WithAnnotation("source", "ci"))
	original := []Entry{{Dependency: "dep", Tags: []string{"existing"}}}
	a.Apply(original)
	if annotatorHasTag(original[0].Tags, "source:ci") {
		t.Error("Apply mutated the original entry")
	}
}

func TestAnnotator_EmptyInput_ReturnsEmpty(t *testing.T) {
	a := NewAnnotator(WithAnnotation("k", "v"))
	out := a.Apply([]Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
