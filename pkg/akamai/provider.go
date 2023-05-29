// Package akamai allows to initialize and set up Akamai Provider
package akamai

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/spf13/cast"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
	"github.com/akamai/terraform-provider-akamai/v4/version"
)

const (
	// ProviderRegistryPath is the path for the provider in the terraform registry
	ProviderRegistryPath = "registry.terraform.io/akamai/akamai"

	// ProviderName is the legacy name of the provider
	// Deprecated: terrform now uses registry paths, the shortest of which would be akamai/akamai"
	ProviderName = "terraform-provider-akamai"
)

// ErrWrongEdgeGridConfiguration is returned when the configuration could not be read
var ErrWrongEdgeGridConfiguration = errors.New("error reading Akamai EdgeGrid configuration")

type (
	provider struct {
		schema.Provider
	}
)

var (
	once sync.Once

	instance *provider
)

// Provider returns the provider function to terraform
func Provider(provs ...subprovider.Subprovider) plugin.ProviderFunc {
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
					"request_limit": {
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc("AKAMAI_REQUEST_LIMIT", 0),
						Type:        schema.TypeInt,
						Description: "The maximum number of API requests to be made per second (0 for no limit)",
					},
				},
				ResourcesMap:       make(map[string]*schema.Resource),
				DataSourcesMap:     make(map[string]*schema.Resource),
				ProviderMetaSchema: make(map[string]*schema.Schema),
			},
		}

		for _, p := range provs {
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
		}

		instance.ConfigureContextFunc = configureProviderContext(&instance.Provider)
	})

	return func() *schema.Provider {
		return &instance.Provider
	}
}

func configureProviderContext(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// generate an operation id so we can correlate all calls to this provider
		opid := uuid.Must(uuid.NewRandom()).String()

		// create a log from the hclog in the context
		log := hclog.FromContext(ctx).With(
			"OperationID", opid,
		)

		cacheEnabled, err := tf.GetBoolValue("cache_enabled", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}
		cache.Enable(cacheEnabled)

		edgercOps := []edgegrid.Option{edgegrid.WithEnv(true)}

		edgercPath, err := tf.GetStringValue("edgerc", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}
		edgercPath = getEdgercPath(edgercPath)

		edgercOps = append(edgercOps, edgegrid.WithFile(edgercPath))
		edgercSection, err := tf.GetStringValue("config_section", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}
		if err == nil {
			edgercOps = append(edgercOps, edgegrid.WithSection(edgercSection))
		}
		envs, err := tf.GetSetValue("config", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}
		if err == nil && len(envs.List()) > 0 {
			envsMap, ok := envs.List()[0].(map[string]interface{})
			if !ok {
				return nil, diag.FromErr(fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "config", "map[string]interface{}"))
			}
			err = setEdgegridEnvs(envsMap, edgercSection)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		requestLimit, err := tf.GetIntValue("request_limit", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}

		edgerc, err := edgegrid.New(edgercOps...)
		if err != nil {
			return nil, diag.Errorf("%s: %s", ErrWrongEdgeGridConfiguration, err.Error())
		}

		if err := edgerc.Validate(); err != nil {
			return nil, diag.Errorf(err.Error())
		}

		// PROVIDER_VERSION env value must be updated in version file, for every new release.
		userAgent := p.UserAgent(ProviderName, version.ProviderVersion)
		logger := logger.FromHCLog(log)
		logger.Infof("Provider version: %s", version.ProviderVersion)

		logger.Debugf("Using request_limit value %d", requestLimit)
		sess, err := session.New(
			session.WithSigner(edgerc),
			session.WithUserAgent(userAgent),
			session.WithLog(logger),
			session.WithHTTPTracing(cast.ToBool(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED"))),
			session.WithRequestLimit(requestLimit),
		)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		meta, err := meta.New(sess, log, opid)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return meta, nil
	}
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
				return fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "access_token", "string")
			}
		case "CLIENT_TOKEN":
			value, ok = envsMap["client_token"].(string)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "client_token", "string")
			}
		case "HOST":
			value, ok = envsMap["host"].(string)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "host", "string")
			}
		case "CLIENT_SECRET":
			value, ok = envsMap["client_secret"].(string)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "client_secret", "string")
			}
		case "MAX_BODY":
			maxBody, ok := envsMap["max_body"].(int)
			if !ok {
				return fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "max_body", "int")
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
