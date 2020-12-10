package iam

import "github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"

// Cache interface needed for this provider
type Cache interface {
	Set(key string, val interface{}) error
	Get(key string, out interface{}) error
}

// Implements Cache as provided by the meta
type metaCache struct {
	sub  akamai.Subprovider
	meta akamai.OperationMeta
}

func (c metaCache) Set(key string, val interface{}) error {
	if c.sub == nil || c.meta == nil {
		return nil
	}

	return c.meta.CacheSet(c.sub, key, val)
}

func (c metaCache) Get(key string, out interface{}) error {
	if c.sub == nil || c.meta == nil {
		return akamai.ErrCacheDisabled
	}

	return c.meta.CacheGet(c.sub, key, out)
}
