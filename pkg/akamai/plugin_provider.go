package akamai

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/collections"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
)

// NewPluginProvider returns the provider function to terraform
func NewPluginProvider(subprovs ...subprovider.Plugin) plugin.ProviderFunc {
	prov := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edgerc": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"config_section": {
				Description: "The section of the edgerc file to use for configuration",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"config": {
				Optional:      true,
				Type:          schema.TypeSet,
				Elem:          config.PluginOptions(),
				MaxItems:      1,
				ConflictsWith: []string{"edgerc", "config_section"},
			},
			"cache_enabled": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"request_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of API requests to be made per second (0 for no limit)",
			},
		},
		ResourcesMap:   make(map[string]*schema.Resource),
		DataSourcesMap: make(map[string]*schema.Resource),
	}

	for _, subprov := range subprovs {
		if err := collections.AddMap(prov.ResourcesMap, subprov.Resources()); err != nil {
			panic(err)
		}

		if err := collections.AddMap(prov.DataSourcesMap, subprov.DataSources()); err != nil {
			panic(err)
		}
	}

	prov.ConfigureContextFunc = configureProviderContext(prov)

	return func() *schema.Provider {
		return prov
	}
}

func configureProviderContext(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		cacheEnabled, err := tf.GetBoolValue("cache_enabled", d)
		if err != nil {
			if !errors.Is(err, tf.ErrNotFound) {
				return nil, diag.FromErr(err)
			}
			cacheEnabled = true
		}
		cache.Enable(cacheEnabled)

		edgercPath, err := tf.GetStringValue("edgerc", d)
		if err != nil {
			if !errors.Is(err, tf.ErrNotFound) {
				return nil, diag.FromErr(err)
			}
			if v := os.Getenv("EDGERC"); v != "" {
				edgercPath = v
			}
		}

		edgercSection, err := tf.GetStringValue("config_section", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}

		envs, err := tf.GetSetValue("config", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}

		var edgercConfig map[string]any
		if err == nil && len(envs.List()) > 0 {
			envsMap, ok := envs.List()[0].(map[string]any)
			if !ok {
				return nil, diag.FromErr(fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "config", "map[string]any"))
			}
			edgercConfig = envsMap
		}

		requestLimit, err := tf.GetIntValue("request_limit", d)
		if err != nil {
			if !errors.Is(err, tf.ErrNotFound) {
				return nil, diag.FromErr(err)
			}
			if v := os.Getenv("AKAMAI_REQUEST_LIMIT"); v != "" {
				requestLimit, err = strconv.Atoi(v)
				if err != nil {
					return nil, diag.FromErr(err)
				}
			}
		}

		meta, err := configureContext(contextConfig{
			edgercPath:    edgercPath,
			edgercSection: edgercSection,
			edgercConfig:  edgercConfig,
			userAgent:     userAgent(p.TerraformVersion),
			ctx:           ctx,
			requestLimit:  requestLimit,
			enableCache:   cacheEnabled,
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return meta, nil
	}
}
