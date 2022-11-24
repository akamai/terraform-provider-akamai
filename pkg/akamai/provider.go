package akamai

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/spf13/cast"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/akamai/terraform-provider-akamai/v2/version"
)

const (
	// ProviderRegistryPath is the path for the provider in the terraform registry
	ProviderRegistryPath = "registry.terraform.io/akamai/akamai"

	// ProviderName is the legacy name of the provider
	// Deprecated: terrform now uses registry paths, the shortest of which would be akamai/akamai"
	ProviderName = "terraform-provider-akamai"

	// ConfigurationIsNotSpecified is the message for when EdgeGrid configuration is not specified
	ConfigurationIsNotSpecified = "Akamai EdgeGrid configuration was not specified. Specify the configuration using system environment variables or the location and file name containing the edgerc configuration. Default location the provider checks for is the current user’s home directory. Default configuration file name the provider checks for is .edgerc."
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
					},
					"config": {
						Optional: true,
						Type:     schema.TypeSet,
						Elem:     config.Options("config"),
						MaxItems: 1,
					},
					"cache_enabled": {
						Optional: true,
						Default:  true,
						Type:     schema.TypeBool,
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
			return configureContext(ctx, d)
		}
	})

	return func() *schema.Provider {
		return &instance.Provider
	}
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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

	cacheEnabled, err := tools.GetBoolValue("cache_enabled", d)
	if err != nil && !IsNotFoundError(err) {
		return nil, diag.FromErr(err)
	}

	edgercOps := []edgegrid.Option{edgegrid.WithEnv(true)}

	edgercPath, err := tools.GetStringValue("edgerc", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, diag.FromErr(err)
	}
	edgercPath = getEdgercPath(edgercPath)

	edgercOps = append(edgercOps, edgegrid.WithFile(edgercPath))
	edgercSection, err := tools.GetStringValue("config_section", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, diag.FromErr(err)
	}
	if err == nil {
		edgercOps = append(edgercOps, edgegrid.WithSection(edgercSection))
	}
	envs, err := tools.GetSetValue("config", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, diag.FromErr(err)
	}
	if err == nil && len(envs.List()) > 0 {
		envsMap, ok := envs.List()[0].(map[string]interface{})
		if !ok {
			return nil, diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "config", "map[string]interface{}"))
		}
		err = setEdgegridEnvs(envsMap, edgercSection)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	edgerc, err := edgegrid.New(edgercOps...)
	if err != nil {
		return nil, diag.Errorf(ConfigurationIsNotSpecified)
	}

	if err := edgerc.Validate(); err != nil {
		return nil, diag.Errorf(err.Error())
	}

	// PROVIDER_VERSION env value must be updated in version file, for every new release.
	userAgent := instance.UserAgent(ProviderName, version.ProviderVersion)
	logger := LogFromHCLog(log)
	logger.Infof("Provider version: %s", version.ProviderVersion)

	sess, err := session.New(
		session.WithSigner(edgerc),
		session.WithUserAgent(userAgent),
		session.WithLog(logger),
		session.WithHTTPTracing(cast.ToBool(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED"))),
	)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	meta := &meta{
		log:          log,
		operationID:  opid,
		sess:         sess,
		cacheEnabled: cacheEnabled,
	}

	return meta, nil
}

func getEdgercPath(edgercPath string) string {
	if edgercPath == "" {
		edgercPath = edgegrid.DefaultConfigFile
	}
	return edgercPath
}

func setEdgegridEnvs(envsMap map[string]interface{}, section string) error {
	configEnvs := []string{"ACCESS_TOKEN", "CLIENT_TOKEN", "HOST", "CLIENT_SECRET", "MAX_BODY"}
	prefix := "AKAMAI"
	if section != "" {
		prefix = fmt.Sprintf("%s_%s", prefix, strings.ToUpper(section))
	}
	for _, env := range configEnvs {
		var value string
		var ok bool
		switch env {
		case "ACCESS_TOKEN":
			value, ok = envsMap["access_token"].(string)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "access_token", "string")
			}
		case "CLIENT_TOKEN":
			value, ok = envsMap["client_token"].(string)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "client_token", "string")
			}
		case "HOST":
			value, ok = envsMap["host"].(string)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "host", "string")
			}
		case "CLIENT_SECRET":
			value, ok = envsMap["client_secret"].(string)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "client_secret", "string")
			}
		case "MAX_BODY":
			maxBody, ok := envsMap["max_body"].(int)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "max_body", "int")
			}
			value = strconv.Itoa(maxBody)
		}
		env = fmt.Sprintf("%s_%s", prefix, env)
		if os.Getenv(env) != "" {
			continue
		}
		if err := os.Setenv(env, value); err != nil {
			return err
		}
	}
	return nil
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
