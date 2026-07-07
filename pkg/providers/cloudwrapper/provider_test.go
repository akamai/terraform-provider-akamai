package cloudwrapper

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

type testSubprovider struct {
	*Subprovider
	activationPollInterval time.Duration
	configPollInterval     time.Duration
	configDeleteTimeout    time.Duration
}

func (s *testSubprovider) FrameworkResources() []func() resource.Resource {
	prodResources := s.Subprovider.FrameworkResources()

	for i, fn := range prodResources {
		prodResources[i] = func() resource.Resource {
			res := fn()
			if r, ok := res.(*activationResource); ok {
				r.activationPollInterval = s.activationPollInterval
			}
			if r, ok := res.(*ConfigurationResource); ok {
				r.pollInterval = s.configPollInterval
				r.deleteTimeout = s.configDeleteTimeout
			}
			return res
		}
	}
	return prodResources
}

func (s *testSubprovider) FrameworkDataSources() []func() datasource.DataSource {
	return s.Subprovider.FrameworkDataSources()
}

func newProviderFactory(client *edgegrid.TestClient, interval ...time.Duration) map[string]func() (tfprotov6.ProviderServer, error) {
	activationInterval := time.Microsecond
	if len(interval) > 0 {
		activationInterval = interval[0]
	}
	sub := &testSubprovider{
		Subprovider:            NewSubprovider(),
		activationPollInterval: activationInterval,
		configPollInterval:     1 * time.Millisecond,
		configDeleteTimeout:    20 * time.Millisecond,
	}
	return testutils.NewTestProtoV6ProviderFactory(client, sub)
}
