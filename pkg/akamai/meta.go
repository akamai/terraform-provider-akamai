package akamai

import (
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/allegro/bigcache/v2"
	"github.com/apex/log"
	"github.com/hashicorp/go-hclog"
)

type (

	// OperationMeta is the akamai meta object interface
	OperationMeta interface {
		// Log constructs an hclog sublogger and returns the log.Interface
		Log(args ...interface{}) log.Interface

		// OperationID returns the operation id
		OperationID() string

		// Session returns the operation API session
		Session() session.Session

		// CacheGet returns an object from the cache
		CacheGet(prov Subprovider, key string, out interface{}) error

		// CacheSet sets a value in the cache
		CacheSet(prov Subprovider, key string, val interface{}) error
	}

	meta struct {
		operationID  string
		log          hclog.Logger
		sess         session.Session
		cacheEnabled bool
	}
)

// Meta return the meta object interface
func Meta(m interface{}) OperationMeta {
	return m.(OperationMeta)
}

// ProviderLog creates a logger for the provider from the meta
func (m *meta) Log(args ...interface{}) log.Interface {
	return LogFromHCLog(m.log.With(args...))
}

// OperationID returns the operation id from the meta
func (m *meta) OperationID() string {
	return m.operationID
}

// Session returns the meta session
func (m *meta) Session() session.Session {
	return m.sess
}

func (m *meta) CacheSet(prov Subprovider, key string, val interface{}) error {
	log := m.Log("meta", "CacheSet")

	if !m.cacheEnabled {
		log.Debug("cache disabled")
		return ErrCacheDisabled
	}

	key = fmt.Sprintf("%s:%s", key, prov.Name())

	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal object to cache: %w", err)
	}

	log.Debugf("cache set for for key %s [%d bytes]", key, len(data))

	return instance.cache.Set(key, data)
}

func (m *meta) CacheGet(prov Subprovider, key string, out interface{}) error {
	log := m.Log("meta", "CacheGet")

	if !m.cacheEnabled {
		log.Debug("cache disabled")
		return ErrCacheDisabled
	}

	key = fmt.Sprintf("%s:%s", key, prov.Name())

	data, err := instance.cache.Get(key)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			log.Debugf("cache miss for for key %s", key)

			return ErrCacheEntryNotFound
		}
		return err
	}

	log.Debugf("cache get for for key %s: [%d bytes]", key, len(data))

	return json.Unmarshal(data, out)
}
