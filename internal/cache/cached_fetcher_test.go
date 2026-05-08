package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/driftcheck/internal/cache"
)

type mockFetcher struct {
	calls  int
	attrs  map[string]interface{}
	err    error
}

func (m *mockFetcher) Fetch(_ context.Context, _, _ string) (map[string]interface{}, error) {
	m.calls++
	return m.attrs, m.err
}

func TestCachedFetcher_CachesResult(t *testing.T) {
	sc := newTempCache(t, time.Minute)
	mock := &mockFetcher{attrs: map[string]interface{}{"ami": "ami-123"}}
	cf := cache.NewCachedFetcher(mock, sc)

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		attrs, err := cf.Fetch(ctx, "aws_instance", "i-abc")
		if err != nil {
			t.Fatalf("Fetch #%d: %v", i, err)
		}
		if attrs["ami"] != "ami-123" {
			t.Errorf("unexpected attrs: %v", attrs)
		}
	}
	if mock.calls != 1 {
		t.Errorf("expected 1 upstream call, got %d", mock.calls)
	}
}

func TestCachedFetcher_PropagatesError(t *testing.T) {
	sc := newTempCache(t, time.Minute)
	mock := &mockFetcher{err: errors.New("api error")}
	cf := cache.NewCachedFetcher(mock, sc)

	_, err := cf.Fetch(context.Background(), "aws_instance", "i-bad")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCachedFetcher_ExpiredRefetches(t *testing.T) {
	sc := newTempCache(t, -time.Second)
	mock := &mockFetcher{attrs: map[string]interface{}{"state": "running"}}
	cf := cache.NewCachedFetcher(mock, sc)

	ctx := context.Background()
	for i := 0; i < 2; i++ {
		_, err := cf.Fetch(ctx, "aws_instance", "i-exp")
		if err != nil {
			t.Fatalf("Fetch: %v", err)
		}
	}
	if mock.calls != 2 {
		t.Errorf("expected 2 upstream calls due to expiry, got %d", mock.calls)
	}
}
