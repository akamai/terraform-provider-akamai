// Package cache contains provider's cache instance
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/akamai/terraform-provider-akamai/v7/pkg/log"
	"github.com/allegro/bigcache/v2"
)

var (
	// ErrDisabled is returned when Get or Set is called on a disabled cache
	ErrDisabled = errors.New("cache disabled")
	// ErrEntryNotFound is returned when object under the given key does not exist
	ErrEntryNotFound = errors.New("cache entry not found")
)

var defaultCache = newCache(10 * time.Minute)

type cache struct {
	cache   *bigcache.BigCache
	enabled bool
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
	defaultCache.enabled = enabled
}

// IsEnabled returns whether cache is enabled
func IsEnabled() bool {
	return defaultCache.enabled
}

// Set sets the given value under the key in cache
func Set(bucket Bucket, key string, val any) error {
	log := log.Get("cache", "CacheSet")

	if !defaultCache.enabled {
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

// Get returns value stored under the key from cache and writes it into out
func Get(bucket Bucket, key string, out any) error {
	log := log.Get("cache", "CacheGet")

	if !defaultCache.enabled {
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

	log.Debugf("cache get for for key %s: [%d bytes]", key, len(data))

	return json.Unmarshal(data, out)
}
