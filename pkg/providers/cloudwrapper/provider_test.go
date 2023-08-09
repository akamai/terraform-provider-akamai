package cloudwrapper

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
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
		resources:   s.Resources(),
		datasources: s.DataSources(),
	}

	for _, opt := range opts {
		opt(ts)
	}

	return ts
}

// Resources returns terraform resources for cloudwrapper
func (ts *TestSubprovider) Resources() []func() resource.Resource {
	for i, fn := range ts.resources {
		// decorate
		fn := fn
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

// DataSources returns terraform data sources for cloudwrapper
func (ts *TestSubprovider) DataSources() []func() datasource.DataSource {
	for i, fn := range ts.datasources {
		fn := fn
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
	if err := testutils.TFTestSetup(); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := testutils.TFTestTeardown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func newProviderFactory(opts ...testSubproviderOption) map[string]func() (tfprotov5.ProviderServer, error) {
	testAccProvider := akamai.NewFrameworkProvider(newTestSubprovider(opts...))()

	return map[string]func() (tfprotov5.ProviderServer, error){
		"akamai": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				providerserver.NewProtocol5(
					testAccProvider,
				),
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
