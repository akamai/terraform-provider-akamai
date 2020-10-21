package akamai

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/allegro/bigcache"
	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

const (
	// ProviderRegistryPath is the path for the provider in the terraform registry
	ProviderRegistryPath = "registry.terraform.io/akamai/akamai"

	// ProviderName is the legacy name of the provider
	// Deprecated: terrform now uses registry paths, the shortest of which would be akamai/akamai"
	ProviderName = "terraform-provider-akamai"
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

		// Configure returns the subprovider opaque state object
		Configure(log.Interface, *schema.ResourceData) diag.Diagnostics
	}

	provider struct {
		schema.Provider
		subs  map[string]Subprovider
		cache *bigcache.BigCache
	}
)

var (
	once sync.Once

	instance *provider
)

// Provider returns the provider function to terraform
func Provider(provs ...Subprovider) plugin.ProviderFunc {
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
			subs: make(map[string]Subprovider),
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
			// generate an operation id so we can correlate all calls to this provider
			opid := uuid.Must(uuid.NewRandom()).String()

			// create a log from the hclog in the context
			log := hclog.FromContext(ctx).With(
				"OperationID", opid,
			)

			// configure sub-providers
			for _, p := range instance.subs {
				if err := p.Configure(LogFromHCLog(log), d); err != nil {
					return nil, err
				}
			}

			edgercOps := []edgegrid.Option{edgegrid.WithEnv(true)}

			edgercPath, err := tools.GetStringValue("edgerc", d)
			if err != nil && !IsNotFoundError(err) {
				return nil, diag.FromErr(err)
			}
			if edgercPath != "" {
				edgercOps = append(edgercOps, edgegrid.WithFile(edgercPath))
			} else {
				edgercOps = append(edgercOps, edgegrid.WithFile(edgegrid.DefaultConfigFile))
			}
			environment, err := tools.GetStringValue("config_section", d)
			if err != nil && !IsNotFoundError(err) {
				return nil, diag.FromErr(err)
			}
			if environment != "" {
				edgercOps = append(edgercOps, edgegrid.WithSection(environment))
			}

			edgercSection, err := tools.GetStringValue("config_section", d)
			if err != nil && !IsNotFoundError(err) {
				return nil, diag.FromErr(err)
			}
			edgercOps = append(edgercOps, edgegrid.WithSection(edgercSection))

			edgerc, err := edgegrid.New(edgercOps...)
			if err != nil {
				return nil, diag.FromErr(fmt.Errorf("failed to load edgegrid config: %w", err))
			}

			userAgent := instance.UserAgent(ProviderName, instance.TerraformVersion)

			sess, err := session.New(
				session.WithSigner(edgerc),
				session.WithUserAgent(userAgent),
				session.WithLog(LogFromHCLog(log)),
				session.WithHTTPTracing(cast.ToBool(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED"))),
			)

			meta := &meta{
				log:         log,
				operationID: opid,
				sess:        sess,
			}

			return meta, nil
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
