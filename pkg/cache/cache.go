// Package cache contains provider's cache instance
package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/akamai/terraform-provider-akamai/v11/pkg/log"
	"github.com/allegro/bigcache/v2"
)

var deletedSentinel = []byte("__deleted__")

var (
	// ErrDisabled is returned when Get or Set is called on a disabled cache
	ErrDisabled = errors.New("cache disabled")
	// ErrEntryNotFound is returned when object under the given key does not exist
	ErrEntryNotFound = errors.New("cache entry not found")
)

var defaultCache = newCache(10 * time.Minute)

type cache struct {
	cache   *bigcache.BigCache
	enabled atomic.Bool
}

// BucketName can be used as a bucket argument to Set and Get functions
type BucketName string

// Name returns BucketID as a string
func (b BucketName) Name() string {
	return string(b)
}

// Bucket defines a contract for a bucket used to form a key
type Bucket interface {
	Name() string
}

func newCache(eviction time.Duration) *cache {
	c, err := bigcache.NewBigCache(bigcache.DefaultConfig(eviction))
	if err != nil {
		panic(err)
	}

	return &cache{cache: c}
}

// Enable is used to enable or disable cache
func Enable(enabled bool) {
	defaultCache.enabled.Store(enabled)
}

// IsEnabled returns whether cache is enabled
func IsEnabled() bool {
	return defaultCache.enabled.Load()
}

// Set sets the given value under the key in cache
func Set(bucket Bucket, key string, val any) error {
	log := log.Get("cache", "CacheSet")

	if !defaultCache.enabled.Load() {
		log.Debug("cache disabled")
		return ErrDisabled
	}

	key = fmt.Sprintf("%s:%s", key, bucket.Name())

	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal object to cache: %w", err)
	}

	log.Debugf("cache set for for key %s [%d bytes]", key, len(data))

	return defaultCache.cache.Set(key, data)
}

// Delete marks a cache entry as deleted. Subsequent Get calls for the same key
// return ErrEntryNotFound until a new Set overwrites the entry.
func Delete(bucket Bucket, key string) error {
	log := log.Get("cache", "CacheDelete")

	if !defaultCache.enabled.Load() {
		log.Debug("cache disabled")
		return ErrDisabled
	}

	key = fmt.Sprintf("%s:%s", key, bucket.Name())
	log.Debugf("cache delete for key %s", key)
	return defaultCache.cache.Set(key, deletedSentinel)
}

// Get returns value stored under the key from cache and writes it into out
func Get(bucket Bucket, key string, out any) error {
	log := log.Get("cache", "CacheGet")

	if !defaultCache.enabled.Load() {
		log.Debug("cache disabled")
		return ErrDisabled
	}

	key = fmt.Sprintf("%s:%s", key, bucket.Name())

	data, err := defaultCache.cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			log.Debugf("cache miss for key %s", key)
			return ErrEntryNotFound
		}
		return err
	}

	if bytes.Equal(data, deletedSentinel) {
		log.Debugf("cache entry deleted for key %s", key)
		return ErrEntryNotFound
	}

	log.Debugf("cache get for for key %s: [%d bytes]", key, len(data))

	return json.Unmarshal(data, out)
}
