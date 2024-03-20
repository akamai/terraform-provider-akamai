package akamai

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/collections"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// NewSDKProvider returns the provider function to terraform
func NewSDKProvider(subprovs ...subprovider.Subprovider) plugin.ProviderFunc {
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
				MaxItems:      1,
				ConflictsWith: []string{"edgerc", "config_section"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_token": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_token": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Required: true,
						},
						"max_body": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"account_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
			"retry_max": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number retires of API requests, default 10",
			},
			"retry_wait_min": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The minimum wait time in seconds between API requests retries, default is 1 sec",
			},
			"retry_wait_max": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum wait time in seconds between API requests retries, default is 30 sec",
			},
			"retry_disabled": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Should the retries of API requests be disabled, default false",
			},
		},
		ResourcesMap:   make(map[string]*schema.Resource),
		DataSourcesMap: make(map[string]*schema.Resource),
	}

	for _, subprov := range subprovs {
		if err := collections.AddMap(prov.ResourcesMap, subprov.SDKResources()); err != nil {
			panic(err)
		}

		if err := collections.AddMap(prov.DataSourcesMap, subprov.SDKDataSources()); err != nil {
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

		configSet, err := tf.GetSetValue("config", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, diag.FromErr(err)
		}

		var edgegridConfigBearer configBearer
		if configSet.Len() > 0 {
			configMap, ok := configSet.List()[0].(map[string]any)
			if !ok {
				return nil, diag.FromErr(fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "config", "map[string]any"))
			}
			edgegridConfigBearer = configBearer{
				accessToken:  configMap["access_token"].(string),
				accountKey:   configMap["account_key"].(string),
				clientSecret: configMap["client_secret"].(string),
				clientToken:  configMap["client_token"].(string),
				host:         configMap["host"].(string),
				maxBody:      configMap["max_body"].(int),
			}
		}

		edgegridConfig, err := newEdgegridConfig(edgercPath, edgercSection, edgegridConfigBearer)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		requestLimit, err := getPluginConfigInt(d, "request_limit", "AKAMAI_REQUEST_LIMIT")
		if err != nil {
			return nil, diag.FromErr(err)
		}

		retryMax, err := getPluginConfigInt(d, "retry_max", "AKAMAI_RETRY_MAX")
		if err != nil {
			return nil, diag.FromErr(err)
		}

		retryWaitMin, err := getPluginConfigInt(d, "retry_wait_min", "AKAMAI_RETRY_WAIT_MIN")
		if err != nil {
			return nil, diag.FromErr(err)
		}

		retryWaitMax, err := getPluginConfigInt(d, "retry_wait_max", "AKAMAI_RETRY_WAIT_MAX")
		if err != nil {
			return nil, diag.FromErr(err)
		}

		retryDisabled, err := getPluginConfigBool(d, "retry_disabled", "AKAMAI_RETRY_DISABLED")
		if err != nil {
			return nil, diag.FromErr(err)
		}

		meta, err := configureContext(contextConfig{
			edgegridConfig: edgegridConfig,
			userAgent:      userAgent(p.TerraformVersion),
			ctx:            ctx,
			requestLimit:   requestLimit,
			enableCache:    cacheEnabled,
			retryMax:       retryMax,
			retryWaitMin:   time.Duration(retryWaitMin) * time.Second,
			retryWaitMax:   time.Duration(retryWaitMax) * time.Second,
			retryDisabled:  retryDisabled,
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return meta, nil
	}
}

// NewProtoV6SDKProvider upgrades SDK provider from protocol version 5 to 6
func NewProtoV6SDKProvider(subproviders []subprovider.Subprovider) (func() tfprotov6.ProviderServer, error) {
	pluginProvider, err := tf5to6server.UpgradeServer(
		context.Background(),
		NewSDKProvider(subproviders...)().GRPCProvider,
	)
	return func() tfprotov6.ProviderServer {
		return pluginProvider
	}, err
}

func getPluginConfigInt(d *schema.ResourceData, key string, envKey string) (int, error) {
	value, err := tf.GetIntValue(key, d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return 0, err
		}
		if v := os.Getenv(envKey); v != "" {
			value, err = strconv.Atoi(v)
			if err != nil {
				return 0, err
			}
		}
	}
	return value, nil
}

func getPluginConfigBool(d *schema.ResourceData, key string, envKey string) (bool, error) {
	value, err := tf.GetBoolValue(key, d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return false, err
		}
		if v := os.Getenv(envKey); v != "" {
			value, err = strconv.ParseBool(v)
			if err != nil {
				return false, err
			}
		}
	}
	return value, nil
}
