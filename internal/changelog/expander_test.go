package changelog

import (
	"testing"
	"time"
)

func TestNewExpander_NotNil(t *testing.T) {
	if NewExpander() == nil {
		t.Fatal("expected non-nil Expander")
	}
}

func TestExpander_FullVersion_Unchanged(t *testing.T) {
	e := NewExpander()
	entries := []Entry{{Version: "v1.2.3"}}
	out := e.Apply(entries)
	if out[0].Version != "v1.2.3" {
		t.Errorf("got %q, want %q", out[0].Version, "v1.2.3")
	}
}

func TestExpander_TwoSegments_PadsZero(t *testing.T) {
	e := NewExpander()
	entries := []Entry{{Version: "v1.2"}}
	out := e.Apply(entries)
	if out[0].Version != "v1.2.0" {
		t.Errorf("got %q, want %q", out[0].Version, "v1.2.0")
	}
}

func TestExpander_OneSegment_PadsTwoZeros(t *testing.T) {
	e := NewExpander()
	entries := []Entry{{Version: "v1"}}
	out := e.Apply(entries)
	if out[0].Version != "v1.0.0" {
		t.Errorf("got %q, want %q", out[0].Version, "v1.0.0")
	}
}

func TestExpander_NoVPrefix(t *testing.T) {
	e := NewExpander()
	entries := []Entry{{Version: "2.1"}}
	out := e.Apply(entries)
	if out[0].Version != "2.1.0" {
		t.Errorf("got %q, want %q", out[0].Version, "2.1.0")
	}
}

func TestExpander_PreReleaseSuffix_Preserved(t *testing.T) {
	e := NewExpander()
	entries := []Entry{{Version: "v1.2-alpha.1"}}
	out := e.Apply(entries)
	if out[0].Version != "v1.2.0-alpha.1" {
		t.Errorf("got %q, want %q", out[0].Version, "v1.2.0-alpha.1")
	}
}

func TestExpander_EmptyVersion_Unchanged(t *testing.T) {
	e := NewExpander()
	entries := []Entry{{Version: ""}}
	out := e.Apply(entries)
	if out[0].Version != "" {
		t.Errorf("expected empty version, got %q", out[0].Version)
	}
}

func TestExpander_MultipleEntries_AllExpanded(t *testing.T) {
	e := NewExpander()
	entries := []Entry{
		{Version: "v1", Date: time.Now()},
		{Version: "v2.3", Date: time.Now()},
		{Version: "v4.5.6", Date: time.Now()},
	}
	out := e.Apply(entries)
	want := []string{"v1.0.0", "v2.3.0", "v4.5.6"}
	for i, w := range want {
		if out[i].Version != w {
			t.Errorf("entry %d: got %q, want %q", i, out[i].Version, w)
		}
	}
}
