package edgegrid

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets/v3"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlskeystore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

var (
	_ Client = &TestClient{}
)

// TestClient is a mock implementation of the Client interface for testing purposes.
type TestClient struct {
	CloudCertificates *cloudcertificates.Mock
	CloudletsV2       *cloudlets.Mock
	CloudletsV3       *v3.Mock
	CPS               *cps.Mock
	DomainOwnership   *domainownership.Mock
	HAPI              *hapi.Mock
	IAM               *iam.Mock
	MTLSKeystore      *mtlskeystore.Mock
	PAPI              *papi.Mock
}

// NewTestClient creates a new instance of TestClient with mock implementations.
func NewTestClient() *TestClient {
	return &TestClient{
		CloudCertificates: &cloudcertificates.Mock{},
		CloudletsV2:       &cloudlets.Mock{},
		CloudletsV3:       &v3.Mock{},
		CPS:               &cps.Mock{},
		DomainOwnership:   &domainownership.Mock{},
		HAPI:              &hapi.Mock{},
		IAM:               &iam.Mock{},
		MTLSKeystore:      &mtlskeystore.Mock{},
		PAPI:              &papi.Mock{},
	}
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

// GetDomainOwnership returns the mock Domain Ownership client.
func (c *TestClient) GetDomainOwnership() domainownership.DomainOwnership {
	return c.DomainOwnership
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

// GetPAPI returns the mock PAPI client.
func (c *TestClient) GetPAPI() papi.PAPI {
	return c.PAPI
}
