// Package cps contains implementation for Akamai Terraform sub-provider responsible for maintaining certificates
package cps

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/subprovider"
)

type (
	// Subprovider gathers CPS resources and data sources
	Subprovider struct{}
)

var _ subprovider.Subprovider = &Subprovider{}

const (
	// Default polling intervals
	defaultPollChangeStatusInterval  = 10 * time.Second
	defaultPollGetEnrollmentInterval = 30 * time.Second
)

// NewSubprovider returns a new CPS subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// SDKResources returns the CPS resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cps_dv_enrollment":          resourceCPSDVEnrollment(defaultPollChangeStatusInterval, defaultPollGetEnrollmentInterval),
		"akamai_cps_dv_validation":          resourceCPSDVValidation(defaultPollChangeStatusInterval),
		"akamai_cps_third_party_enrollment": resourceCPSThirdPartyEnrollment(defaultPollChangeStatusInterval, defaultPollGetEnrollmentInterval),
		"akamai_cps_upload_certificate":     resourceCPSUploadCertificate(defaultPollChangeStatusInterval),
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
