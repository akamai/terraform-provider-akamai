// Package cps contains implementation for Akamai Terraform sub-provider responsible for maintaining certificates
package cps

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
)

type (
	// Subprovider gathers CPS resources and data sources
	Subprovider struct {
		client cps.CPS
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new CPS subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Client returns the CPS interface
func (p *Subprovider) Client(meta meta.Meta) cps.CPS {
	if p.client != nil {
		return p.client
	}
	return cps.Client(meta.Session())
}

// SDKResources returns the CPS resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cps_dv_enrollment":          resourceCPSDVEnrollment(),
		"akamai_cps_dv_validation":          resourceCPSDVValidation(),
		"akamai_cps_third_party_enrollment": resourceCPSThirdPartyEnrollment(),
		"akamai_cps_upload_certificate":     resourceCPSUploadCertificate(),
	}
}

// SDKDataSources returns the CPS data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cps_csr":         dataSourceCPSCSR(),
		"akamai_cps_deployments": dataSourceDeployments(),
		"akamai_cps_enrollment":  dataSourceCPSEnrollment(),
		"akamai_cps_enrollments": dataSourceCPSEnrollments(),
		"akamai_cps_warnings":    dataSourceCPSWarnings(),
	}
}

// FrameworkResources returns the CPS resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the CPS data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
