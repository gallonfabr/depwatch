package changelog

import (
	"testing"
)

func TestParseVersion_Valid(t *testing.T) {
	v, err := ParseVersion("v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("got %v, want 1.2.3", v)
	}
}

func TestParseVersion_NoVPrefix(t *testing.T) {
	v, err := ParseVersion("0.10.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 0 || v.Minor != 10 || v.Patch != 5 {
		t.Errorf("got %v, want 0.10.5", v)
	}
}

func TestParseVersion_PreRelease(t *testing.T) {
	v, err := ParseVersion("v2.0.0-beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Pre != "beta" {
		t.Errorf("expected pre=beta, got %q", v.Pre)
	}
}

func TestParseVersion_Invalid(t *testing.T) {
	cases := []string{"1.2", "abc", "", "1.x.3"}
	for _, tc := range cases {
		_, err := ParseVersion(tc)
		if err == nil {
			t.Errorf("expected error for %q, got nil", tc)
		}
	}
}

func TestVersion_String(t *testing.T) {
	v := Version{Major: 3, Minor: 1, Patch: 4}
	if v.String() != "3.1.4" {
		t.Errorf("got %q, want \"3.1.4\"", v.String())
	}

	v.Pre = "rc1"
	if v.String() != "3.1.4-rc1" {
		t.Errorf("got %q, want \"3.1.4-rc1\"", v.String())
	}
}

func TestVersion_Less(t *testing.T) {
	a, _ := ParseVersion("1.0.0")
	b, _ := ParseVersion("2.0.0")
	if !a.Less(b) {
		t.Error("expected 1.0.0 < 2.0.0")
	}
	if b.Less(a) {
		t.Error("expected 2.0.0 not < 1.0.0")
	}
}

func TestVersion_Equal(t *testing.T) {
	a, _ := ParseVersion("1.2.3")
	b, _ := ParseVersion("v1.2.3-alpha")
	if !a.Equal(b) {
		t.Error("expected versions to be equal ignoring pre-release")
	}
}
