package cache

import (
	"context"
	"fmt"
	"log"
)

// Fetcher is the interface for retrieving live cloud resource attributes.
type Fetcher interface {
	Fetch(ctx context.Context, resourceType, resourceID string) (map[string]interface{}, error)
}

// CachedFetcher wraps a Fetcher and transparently caches results.
type CachedFetcher struct {
	inner Fetcher
	cache *StateCache
}

// NewCachedFetcher creates a CachedFetcher backed by the given Fetcher and StateCache.
func NewCachedFetcher(inner Fetcher, sc *StateCache) *CachedFetcher {
	return &CachedFetcher{inner: inner, cache: sc}
}

// Fetch returns cached attributes when available, otherwise delegates to the
// underlying Fetcher and stores the result.
func (cf *CachedFetcher) Fetch(ctx context.Context, resourceType, resourceID string) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("%s.%s", resourceType, resourceID)

	if entry, ok := cf.cache.Get(cacheKey); ok {
		log.Printf("cache hit: %s", cacheKey)
		return entry.Attributes, nil
	}

	attrs, err := cf.inner.Fetch(ctx, resourceType, resourceID)
	if err != nil {
		return nil, err
	}

	entry := &Entry{
		ResourceID: cacheKey,
		Attributes: attrs,
	}
	if setErr := cf.cache.Set(entry); setErr != nil {
		log.Printf("cache: failed to store %s: %v", cacheKey, setErr)
	}
	return attrs, nil
}
