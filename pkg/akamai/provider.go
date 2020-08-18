package akamai

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/allegro/bigcache"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mr-tron/base58"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

type (
	// Subprovider is the interface implemented by the sub providers
	Subprovider interface {
		// Name should return the name of the subprovider
		Name() string

		// Version returns the version of the subprovider
		Version() string

		// Schema returns the schemas for the subprovider
		Schema() map[string]*schema.Schema

		// Resources returns the resources for the subprovider
		Resources() map[string]*schema.Resource

		// DataSources returns the datasources for the subprovider
		DataSources() map[string]*schema.Resource

		// Configure returns the subprovider opqaue state object
		Configure(context.Context, hclog.Logger, *schema.ResourceData) (interface{}, diag.Diagnostics)
	}

	// Context provides logging and other support services to the adapters
	Context interface {
		// Log returns a named logger for the subprovider
		Log(prefix ...string) hclog.Logger

		// Meta returns this providers internal meta object
		Meta() interface{}

		// CacheSet sets an object in the meta cache
		CacheSet(key string, value interface{}) error

		// CacheGet gets an object from the meta cache
		CacheGet(key string, out interface{}) error

		// OperationID is a unique id for an operation
		OperationID() string

		// TerraformVersion returns the version from the core provider
		TerraformVersion() string
	}

	provider struct {
		schema.Provider
		log    hclog.Logger
		subs   map[string]Subprovider
		states map[string]interface{}
		cache  *bigcache.BigCache
	}

	akaContext struct {
		operationID string
		log         hclog.Logger
		meta        interface{}
	}
)

var (
	once sync.Once

	instance *provider
)

// Provider returns the provider function to terraform
func Provider(log hclog.Logger, provs ...Subprovider) plugin.ProviderFunc {
	once.Do(func() {
		instance = &provider{
			Provider: schema.Provider{
				Schema: map[string]*schema.Schema{
					"edgerc": {
						Optional:    true,
						Type:        schema.TypeString,
						DefaultFunc: schema.EnvDefaultFunc("EDGERC", nil),
					},
					"config_section": {
						Description: "The section of the edgerc file to use for configuration",
						Optional:    true,
						Type:        schema.TypeString,
						Default:     "default",
					},
				},
				ResourcesMap:       make(map[string]*schema.Resource),
				DataSourcesMap:     make(map[string]*schema.Resource),
				ProviderMetaSchema: make(map[string]*schema.Schema),
			},
			subs:   make(map[string]Subprovider),
			states: make(map[string]interface{}),
			log:    log,
		}

		cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
		if err != nil {
			panic(err)
		}

		instance.cache = cache

		for _, p := range provs {
			subSchema, err := mergeSchema(p.Schema(), instance.Schema)
			if err != nil {
				panic(err)
			}
			instance.Schema = subSchema
			resources, err := mergeResource(p.Resources(), instance.ResourcesMap)
			if err != nil {
				panic(err)
			}
			instance.ResourcesMap = resources
			dataSources, err := mergeResource(p.DataSources(), instance.DataSourcesMap)
			if err != nil {
				panic(err)
			}
			instance.DataSourcesMap = dataSources

			instance.subs[p.Name()] = p
		}

		instance.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			var stateSet bool

			for _, p := range instance.subs {
				state, err := p.Configure(ctx, log, d)
				if err != nil {
					return nil, err
				}

				if state != nil {
					stateSet = true
					instance.states[p.Name()] = state
				}
			}

			if !stateSet {
				return nil, ErrNoConfiguredProviders.Diagnostics()
			}

			// TODO: once the client is update this will be done elsewhere
			client.UserAgent = instance.UserAgent("terraform-provider-akamai", instance.TerraformVersion)

			return &instance, nil
		}
	})

	return func() *schema.Provider {
		return &instance.Provider
	}
}

func mergeSchema(from, to map[string]*schema.Schema) (map[string]*schema.Schema, error) {
	for k, v := range from {
		if _, ok := to[k]; ok {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateSchemaKey, k)
		}
		to[k] = v
	}
	return to, nil
}

func mergeResource(from, to map[string]*schema.Resource) (map[string]*schema.Resource, error) {
	for k, v := range from {
		if _, ok := to[k]; ok {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateSchemaKey, k)
		}
		to[k] = v
	}
	return to, nil
}

// ContextGet returns the context object from the passed interface
func ContextGet(name string) (Context, bool) {
	sub, ok := instance.subs[name]
	if !ok {
		return nil, false
	}

	coid := uuid.Must(uuid.NewRandom())
	opid := base58.Encode(coid[:])
	m := akaContext{
		operationID: opid,
		log:         instance.log.Named(sub.Name()).Named(sub.Version()).Named(opid),
	}

	if state, ok := instance.states[name]; ok {
		m.meta = state
	}

	return &m, true
}

func (c *akaContext) Log(prefix ...string) hclog.Logger {
	if len(prefix) > 0 {
		log := c.log
		for _, p := range prefix {
			log = log.Named(p)
		}
		return log
	}
	return c.log
}

func (c *akaContext) Meta() interface{} {
	return c.meta
}

func (c *akaContext) OperationID() string {
	return c.operationID
}

func (c *akaContext) TerraformVersion() string {
	return instance.TerraformVersion
}

func (c *akaContext) CacheSet(key string, val interface{}) error {
	var in []byte

	switch v := val.(type) {
	case []byte:
		in = v
	default:
		data, err := json.Marshal(val)
		if err != nil {
			return err
		}

		in = data
	}

	return instance.cache.Set(key, in)
}

func (c *akaContext) CacheGet(key string, out interface{}) error {
	data, err := instance.cache.Get(key)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return ErrCacheEntryNotFound(key)
		}

		return err
	}

	return json.Unmarshal(data, out)
}
