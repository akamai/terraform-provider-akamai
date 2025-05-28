package cloudwrapper

import (
	"context"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	TestSubprovider struct {
		resources   []func() resource.Resource
		datasources []func() datasource.DataSource
		client      cloudwrapper.CloudWrapper
		interval    time.Duration
	}

	clientSetter interface {
		setClient(cloudwrapper.CloudWrapper)
	}

	pollIntervalSetter interface {
		setPollInterval(time.Duration)
	}

	testSubproviderOption func(*TestSubprovider)
)

func withMockClient(mock cloudwrapper.CloudWrapper) testSubproviderOption {
	return func(ts *TestSubprovider) {
		ts.client = mock
		ts.interval = time.Microsecond
	}
}

func withInterval(interval time.Duration) testSubproviderOption {
	return func(ts *TestSubprovider) {
		ts.interval = interval
	}
}

func newTestSubprovider(opts ...testSubproviderOption) *TestSubprovider {
	s := NewSubprovider()

	ts := &TestSubprovider{
		resources:   s.FrameworkResources(),
		datasources: s.FrameworkDataSources(),
	}

	for _, opt := range opts {
		opt(ts)
	}

	return ts
}

func (ts *TestSubprovider) SDKResources() map[string]*schema.Resource {
	return nil
}

func (ts *TestSubprovider) SDKDataSources() map[string]*schema.Resource {
	return nil
}

func (ts *TestSubprovider) FrameworkResources() []func() resource.Resource {
	for i, fn := range ts.resources {
		// decorate
		ts.resources[i] = func() resource.Resource {
			res := fn()
			if v, ok := res.(clientSetter); ok {
				v.setClient(ts.client)
			}
			if v, ok := res.(pollIntervalSetter); ok {
				v.setPollInterval(ts.interval)
			}
			return res
		}
	}
	return ts.resources
}

func (ts *TestSubprovider) FrameworkDataSources() []func() datasource.DataSource {
	for i, fn := range ts.datasources {
		// decorate
		ts.datasources[i] = func() datasource.DataSource {
			ds := fn()
			if v, ok := ds.(clientSetter); ok {
				v.setClient(ts.client)
			}
			return ds
		}
	}
	return ts.datasources
}

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

func newProviderFactory(opts ...testSubproviderOption) map[string]func() (tfprotov6.ProviderServer, error) {
	testAccProvider := akamai.NewFrameworkProvider(newTestSubprovider(opts...))()

	return map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(
					testAccProvider,
				),
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
