package mtlskeystore

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

type (
	// CustomPollingSubprovider is an MTLSKeystore subprovider with customizable polling parameters
	// for testing purposes.
	CustomPollingSubprovider struct {
		Subprovider
		pollingInterval time.Duration
		pollingTimeout  time.Duration
	}
)

func NewCustomPollingSubprovider(pollingInterval, pollingTimeout time.Duration) *CustomPollingSubprovider {
	return &CustomPollingSubprovider{
		pollingInterval: pollingInterval,
		pollingTimeout:  pollingTimeout,
	}
}

func (p *CustomPollingSubprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewClientCertificateAkamaiResource,
		NewClientCertificateThirdPartyResource,
		NewClientCertificateUploadResourceCustomPolling(p.pollingInterval, p.pollingTimeout),
	}
}

func NewClientCertificateUploadResourceCustomPolling(pollingInterval, pollingTimeout time.Duration) func() resource.Resource {
	return func() resource.Resource {
		return &clientCertificateUploadResource{
			pollingInterval: pollingInterval,
			pollingTimeout:  pollingTimeout,
		}
	}
}
