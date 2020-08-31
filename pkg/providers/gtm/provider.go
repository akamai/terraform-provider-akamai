package gtm

import (
	"context"
	"errors"
	"fmt"
	"sync"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider
	}
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider() akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"gtm_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("gtm_section"),
			},
			"gtm": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("gtm"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_gtm_domain":     resourceGTMv1Domain(),
			"akamai_gtm_datacenter": resourceGTMv1Datacenter(),
			"akamai_gtm_property":   resourceGTMv1Property(),
			"akamai_gtm_resource":   resourceGTMv1Resource(),
			"akamai_gtm_cidrmap":    resourceGTMv1Cidrmap(),
			"akamai_gtm_geomap":     resourceGTMv1Geomap(),
			"akamai_gtm_asmap":      resourceGTMv1ASmap(),
		},
	}

	return provider
}

type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
}

func getConfigGTMV1Service(d resourceData) (*edgegrid.Config, error) {
	var GTMv1Config edgegrid.Config
	var err error
	if _, ok := d.GetOk("gtm"); ok {
		resourceGtm, ok := d.Get("gtm").(set)
		if !ok {
			return nil, fmt.Errorf("wrong cast")
		}

		resourceConfig, ok := resourceGtm.List()[0].(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("wrong cast")
		}

		host, ok := resourceConfig["host"].(string)

		if !ok {
			return nil, fmt.Errorf("wrong host data type")
		}

		accessToken, ok := resourceConfig["access_token"].(string)

		if !ok {
			return nil, fmt.Errorf("wrong access_token data type")
		}

		clientToken, ok := resourceConfig["client_token"].(string)

		if !ok {
			return nil, fmt.Errorf("wrong client_token data type")
		}

		clientSecret, ok := resourceConfig["client_secret"].(string)

		if !ok {
			return nil, fmt.Errorf("wrong client_secret data type")
		}

		maxBody, ok := resourceConfig["max_body"].(int)

		if !ok {
			return nil, fmt.Errorf("wrong max_body data type")
		}

		GTMv1Config = edgegrid.Config{
			Host:         host,
			AccessToken:  accessToken,
			ClientToken:  clientToken,
			ClientSecret: clientSecret,
			MaxBody:      maxBody,
		}

		gtm.Init(GTMv1Config)
		edgegrid.SetupLogging()
		return &GTMv1Config, nil
	}

	edgerc, err := tools.GetStringValue("edgerc", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}

	var section string

	for _, s := range tools.FindStringValues(d, "gtm_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	GTMv1Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	gtm.Init(GTMv1Config)
	return &GTMv1Config, nil
}

func (p *provider) Name() string {
	return "gtm"
}

// GtmProviderVersion update version string anytime provider adds new features
const GtmProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return GtmProviderVersion
}

func (p *provider) Schema() map[string]*schema.Schema {
	return p.Provider.Schema
}

func (p *provider) Resources() map[string]*schema.Resource {
	return p.Provider.ResourcesMap
}

func (p *provider) DataSources() map[string]*schema.Resource {
	return p.Provider.DataSourcesMap
}

func (p *provider) Configure(ctx context.Context, log hclog.Logger, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Named(p.Name()).Debug("START Configure")

	cfg, err := getConfigGTMV1Service(d)
	if err != nil {
		return nil, nil
	}

	return cfg, nil
}
