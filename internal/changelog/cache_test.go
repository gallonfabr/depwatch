package changelog

import (
	"testing"
	"time"
)

func TestNewCache_DefaultTTL(t *testing.T) {
	c := NewCache(0)
	if c.TTL != 10*time.Minute {
		t.Errorf("expected default TTL 10m, got %v", c.TTL)
	}
}

func TestNewCache_CustomTTL(t *testing.T) {
	c := NewCache(5 * time.Minute)
	if c.TTL != 5*time.Minute {
		t.Errorf("expected TTL 5m, got %v", c.TTL)
	}
}

func TestCache_SetAndGet_Hit(t *testing.T) {
	c := NewCache(1 * time.Minute)
	c.Set("https://example.com/changelog", "## v1.0.0")

	content, ok := c.Get("https://example.com/changelog")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if content != "## v1.0.0" {
		t.Errorf("unexpected content: %q", content)
	}
}

func TestCache_Get_Miss(t *testing.T) {
	c := NewCache(1 * time.Minute)

	_, ok := c.Get("https://example.com/missing")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestCache_Get_Expired(t *testing.T) {
	c := NewCache(1 * time.Millisecond)
	c.Set("key", "value")

	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := NewCache(1 * time.Minute)
	c.Set("key", "value")
	c.Invalidate("key")

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected invalidated entry to be a miss")
	}
}

func TestCache_Len(t *testing.T) {
	c := NewCache(1 * time.Minute)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	c.Set("a", "1")
	c.Set("b", "2")
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestCache_Overwrite(t *testing.T) {
	c := NewCache(1 * time.Minute)
	c.Set("key", "old")
	c.Set("key", "new")

	content, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if content != "new" {
		t.Errorf("expected 'new', got %q", content)
	}
	if c.Len() != 1 {
		t.Errorf("expected len 1, got %d", c.Len())
	}
}
