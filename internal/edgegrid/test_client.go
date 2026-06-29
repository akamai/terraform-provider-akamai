package edgegrid

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/clientlists"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudlets/v3"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/domainownership"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlskeystore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
)

var (
	_ Client = &TestClient{}
)

// TestClient is a mock implementation of the Client interface for testing purposes.
type TestClient struct {
	ClientLists       *clientlists.Mock
	CloudCertificates *cloudcertificates.Mock
	CloudletsV2       *cloudlets.Mock
	CloudletsV3       *v3.Mock
	CPS               *cps.Mock
	DNS               *dns.Mock
	DomainOwnership   *domainownership.Mock
	EdgeWorkers       *edgeworkers.Mock
	GTM               *gtm.Mock
	HAPI              *hapi.Mock
	IAM               *iam.Mock
	MTLSKeystore      *mtlskeystore.Mock
	MTLSTruststore    *mtlstruststore.Mock
	PAPI              *papi.Mock
	ReportingGroups   *reportinggroups.Mock
}

// NewTestClient creates a new instance of TestClient with mock implementations.
func NewTestClient() *TestClient {
	return &TestClient{
		ClientLists:       &clientlists.Mock{},
		CloudCertificates: &cloudcertificates.Mock{},
		CloudletsV2:       &cloudlets.Mock{},
		CloudletsV3:       &v3.Mock{},
		CPS:               &cps.Mock{},
		DNS:               &dns.Mock{},
		DomainOwnership:   &domainownership.Mock{},
		EdgeWorkers:       &edgeworkers.Mock{},
		GTM:               &gtm.Mock{},
		HAPI:              &hapi.Mock{},
		IAM:               &iam.Mock{},
		MTLSKeystore:      &mtlskeystore.Mock{},
		MTLSTruststore:    &mtlstruststore.Mock{},
		PAPI:              &papi.Mock{},
		ReportingGroups:   &reportinggroups.Mock{},
	}
}

// GetClientLists returns the mock Client Lists client.
func (c *TestClient) GetClientLists() clientlists.ClientLists {
	return c.ClientLists
}

// GetCloudCertificates returns the mock CCM client.
func (c *TestClient) GetCloudCertificates() cloudcertificates.CloudCertificates {
	return c.CloudCertificates
}

// GetCloudletsV2 returns the mock Cloudlets V2 client.
func (c *TestClient) GetCloudletsV2() cloudlets.Cloudlets {
	return c.CloudletsV2
}

// GetCloudletsV3 returns the mock Cloudlets V3 client.
func (c *TestClient) GetCloudletsV3() v3.Cloudlets {
	return c.CloudletsV3
}

// GetCPS returns the mock CPS client.
func (c *TestClient) GetCPS() cps.CPS {
	return c.CPS
}

// GetDNS returns the mock DNS client.
func (c *TestClient) GetDNS() dns.DNS {
	return c.DNS
}

// GetDomainOwnership returns the mock Domain Ownership client.
func (c *TestClient) GetDomainOwnership() domainownership.DomainOwnership {
	return c.DomainOwnership
}

// GetEdgeWorkers returns the mock EdgeWorkers client.
func (c *TestClient) GetEdgeWorkers() edgeworkers.Edgeworkers {
	return c.EdgeWorkers
}

// GetGTM returns the mock GTM client.
func (c *TestClient) GetGTM() gtm.GTM {
	return c.GTM
}

// GetHAPI returns the mock HAPI client.
func (c *TestClient) GetHAPI() hapi.HAPI {
	return c.HAPI
}

// GetIAM returns the mock IAM client.
func (c *TestClient) GetIAM() iam.IAM {
	return c.IAM
}

// GetMTLSKeystore returns the mock MTLS Keystore client.
func (c *TestClient) GetMTLSKeystore() mtlskeystore.MTLSKeystore {
	return c.MTLSKeystore
}

// GetMTLSTruststore returns the mock MTLS Truststore client.
func (c *TestClient) GetMTLSTruststore() mtlstruststore.MTLSTruststore {
	return c.MTLSTruststore
}

// GetPAPI returns the mock PAPI client.
func (c *TestClient) GetPAPI() papi.PAPI {
	return c.PAPI
}

// GetReportingGroups returns the mock Reporting Groups client.
func (c *TestClient) GetReportingGroups() reportinggroups.ReportingGroups {
	return c.ReportingGroups
}
