// package tcache implements TokenSet using bigcache
package tcache

import (
	"time"

	"github.com/allegro/bigcache"
)

type (
	TokenCache struct {
		cache *bigcache.BigCache
	}
)

// New returns a TokenCache which implements TokenSet and can be used
// for in-memory token persistance.
func New() *TokenCache {
	tc := &TokenCache{}
	tc.cache, _ = bigcache.NewBigCache(bigcache.Config{
		Shards:             1024,
		LifeWindow:         10 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       1024,
		Verbose:            false,
		HardMaxCacheSize:   1024,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	})
	return tc
}

// Contains implements proxy.TokenSet
func (tc *TokenCache) Contains(k string) (bool, error) {
	val, _ := tc.cache.Get(k)
	return val != nil, nil
}

// Add includes the given token to the set
func (tc *TokenCache) Add(k string, _ time.Time) error {
	return tc.cache.Set(k, []byte(k))
}
