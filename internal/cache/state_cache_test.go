package cache_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/driftcheck/internal/cache"
)

func newTempCache(t *testing.T, ttl time.Duration) *cache.StateCache {
	t.Helper()
	dir, err := os.MkdirTemp("", "driftcheck-cache-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	c, err := cache.New(dir, ttl)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestCache_SetAndGet(t *testing.T) {
	c := newTempCache(t, time.Minute)
	e := &cache.Entry{
		ResourceID: "aws_instance.web",
		Attributes: map[string]interface{}{"instance_type": "t3.micro"},
	}
	if err := c.Set(e); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, ok := c.Get("aws_instance.web")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got.Attributes["instance_type"] != "t3.micro" {
		t.Errorf("unexpected attribute: %v", got.Attributes)
	}
}

func TestCache_Miss(t *testing.T) {
	c := newTempCache(t, time.Minute)
	_, ok := c.Get("aws_instance.missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestCache_Expired(t *testing.T) {
	c := newTempCache(t, -time.Second) // already expired
	e := &cache.Entry{
		ResourceID: "aws_instance.old",
		Attributes: map[string]interface{}{},
	}
	_ = c.Set(e)
	_, ok := c.Get("aws_instance.old")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := newTempCache(t, time.Minute)
	e := &cache.Entry{ResourceID: "aws_instance.web", Attributes: map[string]interface{}{}}
	_ = c.Set(e)
	if err := c.Invalidate("aws_instance.web"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	_, ok := c.Get("aws_instance.web")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestCache_InvalidateNonExistent(t *testing.T) {
	c := newTempCache(t, time.Minute)
	if err := c.Invalidate("aws_instance.ghost"); err != nil {
		t.Fatalf("Invalidate non-existent should not error: %v", err)
	}
}

func TestCache_OverwriteEntry(t *testing.T) {
	c := newTempCache(t, time.Minute)
	first := &cache.Entry{
		ResourceID: "aws_instance.web",
		Attributes: map[string]interface{}{"instance_type": "t3.micro"},
	}
	if err := c.Set(first); err != nil {
		t.Fatalf("Set first: %v", err)
	}
	second := &cache.Entry{
		ResourceID: "aws_instance.web",
		Attributes: map[string]interface{}{"instance_type": "t3.large"},
	}
	if err := c.Set(second); err != nil {
		t.Fatalf("Set second: %v", err)
	}
	got, ok := c.Get("aws_instance.web")
	if !ok {
		t.Fatal("expected cache hit after overwrite")
	}
	if got.Attributes["instance_type"] != "t3.large" {
		t.Errorf("expected overwritten value, got: %v", got.Attributes)
	}
}
