package cloudcertificates

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

type (
	// CustomPollingSubprovider is a CloudCertificates subprovider with customizable polling timeout
	// for testing purposes.
	CustomPollingSubprovider struct {
		Subprovider
		pollingTimeout time.Duration
	}
)

func NewCustomPollingSubprovider(pollingTimeout time.Duration) *CustomPollingSubprovider {
	return &CustomPollingSubprovider{
		pollingTimeout: pollingTimeout,
	}
}

func (p *CustomPollingSubprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewCertificateResource,
		NewUploadSignedCertificateResourceCustomPolling(p.pollingTimeout),
	}
}

func NewUploadSignedCertificateResourceCustomPolling(pollingTimeout time.Duration) func() resource.Resource {
	return func() resource.Resource {
		return &uploadSignedCertificateResource{
			pollingInterval: 1 * time.Millisecond,
			pollingTimeout:  pollingTimeout,
		}
	}
}
