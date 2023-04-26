// Package gtm contains implementation for Akamai Terraform sub-provider responsible for managing Global Traffic Management (GTM) domain configuration and administration
package gtm

import (
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/config"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client gtm.GTM
	}

	// Option is a gtm provider option
	Option func(p *provider)
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
				Optional:   true,
				Type:       schema.TypeSet,
				Elem:       config.Options("gtm"),
				MaxItems:   1,
				Deprecated: akamai.NoticeDeprecatedUseAlias("gtm"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_gtm_datacenter":         dataSourceGTMDatacenter(),
			"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_gtm_domain":     resourceGTMv1Domain(),
			"akamai_gtm_property":   resourceGTMv1Property(),
			"akamai_gtm_datacenter": resourceGTMv1Datacenter(),
			"akamai_gtm_resource":   resourceGTMv1Resource(),
			"akamai_gtm_asmap":      resourceGTMv1ASmap(),
			"akamai_gtm_geomap":     resourceGTMv1Geomap(),
			"akamai_gtm_cidrmap":    resourceGTMv1Cidrmap(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c gtm.GTM) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *provider) Client(meta akamai.OperationMeta) gtm.GTM {
	if p.client != nil {
		return p.client
	}
	return gtm.Client(meta.Session())
}

func getConfigGTMV1Service(d *schema.ResourceData) error {
	var inlineConfig *schema.Set
	for _, key := range []string{"gtm", "config"} {
		opt, err := tf.GetSetValue(key, d)
		if err != nil {
			if !errors.Is(err, tf.ErrNotFound) {
				return err
			}
			continue
		}
		if inlineConfig != nil {
			return fmt.Errorf("only one inline config section can be defined")
		}
		inlineConfig = opt
	}
	if err := d.Set("config", inlineConfig); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}

	for _, s := range tf.FindStringValues(d, "gtm_section", "config_section") {
		if s != "default" && s != "" {
			if err := d.Set("config_section", s); err != nil {
				return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
			}
			break
		}
	}

	return nil
}

func (p *provider) Name() string {
	return "gtm"
}

// GTMProviderVersion update version string anytime provider adds new features
const GTMProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return GTMProviderVersion
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

func (p *provider) Configure(log log.Interface, d *schema.ResourceData) diag.Diagnostics {
	log.Debug("START Configure")

	if err := getConfigGTMV1Service(d); err != nil {
		return nil
	}
	return nil
}
