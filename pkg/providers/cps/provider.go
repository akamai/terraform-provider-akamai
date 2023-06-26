// Package cps contains implementation for Akamai Terraform sub-provider responsible for maintaining certificates
package cps

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
)

type (
	// Subprovider gathers cps resources and data sources
	Subprovider struct {
		client cps.CPS
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Plugin = &Subprovider{}

// NewSubprovider returns a core sub provider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

func withClient(c cps.CPS) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the CPS interface
func (p *Subprovider) Client(meta meta.Meta) cps.CPS {
	if p.client != nil {
		return p.client
	}
	return cps.Client(meta.Session())
}

// Resources returns terraform resources for cps
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cps_dv_enrollment":          resourceCPSDVEnrollment(),
		"akamai_cps_dv_validation":          resourceCPSDVValidation(),
		"akamai_cps_third_party_enrollment": resourceCPSThirdPartyEnrollment(),
		"akamai_cps_upload_certificate":     resourceCPSUploadCertificate(),
	}
}

// DataSources returns terraform data sources for cps
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cps_csr":         dataSourceCPSCSR(),
		"akamai_cps_deployments": dataSourceDeployments(),
		"akamai_cps_enrollment":  dataSourceCPSEnrollment(),
		"akamai_cps_enrollments": dataSourceCPSEnrollments(),
		"akamai_cps_warnings":    dataSourceCPSWarnings(),
	}
}
