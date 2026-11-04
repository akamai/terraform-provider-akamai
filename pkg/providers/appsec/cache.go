package appsec

import (
	"errors"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/cache"
)

// getCached implements the double-checked locking cache pattern for getRatePolicy,
// getMatchTarget, and getReputationProfile.
func getCached[List any, Item any](
	mu *sync.Mutex,
	cacheKey string,
	logger log.Interface,
	onDisabled func() (*Item, error),
	fetch func() (*List, error),
	filter func(*List) (*Item, error),
) (*Item, error) {
	list := new(List)

	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, list)
	if errors.Is(err, cache.ErrDisabled) {
		return onDisabled()
	}
	if err == nil {
		return filter(list)
	}
	if !errors.Is(err, cache.ErrEntryNotFound) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	// The first Get skips the lock when data is already cached.
	// The second Get skips a redundant fetch when another goroutine cached it while you were waiting for the lock.
	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, list)
	if err == nil {
		return filter(list)
	}
	if errors.Is(err, cache.ErrDisabled) {
		return onDisabled()
	}
	if !errors.Is(err, cache.ErrEntryNotFound) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	list, err = fetch()
	if err != nil {
		return nil, err
	}

	if err = cache.Set(cache.BucketName(SubproviderName), cacheKey, list); err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching into cache: %s", err.Error())
		return nil, err
	}

	return filter(list)
}
