// Package edgegrid provides a client for interacting with Akamai's Open API.
package edgegrid

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudlets/v3"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/domainownership"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlskeystore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

// Client is the interface for the Akamai Edgegrid client.
type Client interface {
	GetCloudCertificates() cloudcertificates.CloudCertificates

	GetCloudletsV2() cloudlets.Cloudlets

	GetCloudletsV3() v3.Cloudlets

	GetCPS() cps.CPS

	GetDomainOwnership() domainownership.DomainOwnership

	GetGTM() gtm.GTM

	GetHAPI() hapi.HAPI

	GetIAM() iam.IAM

	GetMTLSKeystore() mtlskeystore.MTLSKeystore

	GetMTLSTruststore() mtlstruststore.MTLSTruststore

	GetPAPI() papi.PAPI

	GetReportingGroups() reportinggroups.ReportingGroups
}

var (
	_ Client = &ClientImpl{}
)

// ClientImpl is the implementation of the Client interface used in production.
type ClientImpl struct {
	sess session.Session
}

// NewClientImpl creates a new instance of ClientImpl with the provided session.
func NewClientImpl(sess session.Session) *ClientImpl {
	return &ClientImpl{
		sess: sess,
	}
}

// GetCloudCertificates returns the CCM client for managing cloud certificates.
func (c *ClientImpl) GetCloudCertificates() cloudcertificates.CloudCertificates {
	return cloudcertificates.Client(c.sess)
}

// GetCloudletsV2 returns the Cloudlets V2 client for managing cloudlets.
func (c *ClientImpl) GetCloudletsV2() cloudlets.Cloudlets {
	return cloudlets.Client(c.sess)
}

// GetCloudletsV3 returns the Cloudlets V3 client for managing cloudlets.
func (c *ClientImpl) GetCloudletsV3() v3.Cloudlets {
	return v3.Client(c.sess)
}

// GetCPS returns the CPS client for managing certificates.
func (c *ClientImpl) GetCPS() cps.CPS {
	return cps.Client(c.sess)
}

// GetDomainOwnership returns the Domain Ownership client for managing domain ownership.
func (c *ClientImpl) GetDomainOwnership() domainownership.DomainOwnership {
	return domainownership.Client(c.sess)
}

// GetGTM returns the GTM client for managing Global Traffic Management.
func (c *ClientImpl) GetGTM() gtm.GTM {
	return gtm.Client(c.sess)
}

// GetHAPI returns the HAPI client for managing hostnames APIs.
func (c *ClientImpl) GetHAPI() hapi.HAPI {
	return hapi.Client(c.sess)
}

// GetIAM returns the IAM client for managing identity and access management.
func (c *ClientImpl) GetIAM() iam.IAM {
	return iam.Client(c.sess)
}

// GetMTLSKeystore returns the MTLS Keystore client for managing mTLS keystores.
func (c *ClientImpl) GetMTLSKeystore() mtlskeystore.MTLSKeystore {
	return mtlskeystore.Client(c.sess)
}

// GetMTLSTruststore returns the MTLS Truststore client for managing mTLS truststores.
func (c *ClientImpl) GetMTLSTruststore() mtlstruststore.MTLSTruststore {
	return mtlstruststore.Client(c.sess)
}

// GetPAPI returns the PAPI client for managing property APIs.
func (c *ClientImpl) GetPAPI() papi.PAPI {
	return papi.Client(c.sess)
}

// GetReportingGroups returns the Reporting Groups client for managing reporting groups.
func (c *ClientImpl) GetReportingGroups() reportinggroups.ReportingGroups {
	return reportinggroups.Client(c.sess)
}
