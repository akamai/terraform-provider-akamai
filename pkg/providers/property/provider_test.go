package property

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProviders         map[string]func() (tfprotov6.ProviderServer, error)
	testAccPluginProvider    *schema.Provider
	testAccFrameworkProvider provider.Provider
)

func TestMain(m *testing.M) {
	testAccPluginProvider = akamai.NewPluginProvider(NewPluginSubprovider())()
	testAccFrameworkProvider = akamai.NewFrameworkProvider(NewFrameworkSubprovider())()

	testAccProviders = map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkProvider, err := tf5to6server.UpgradeServer(
				context.Background(),
				testAccPluginProvider.GRPCProvider,
			)
			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return upgradedSdkProvider
				},
				providerserver.NewProtocol6(
					testAccFrameworkProvider,
				),
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}

	if err := testutils.TFTestSetup(); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := testutils.TFTestTeardown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(papiCli papi.PAPI, hapiCli hapi.HAPI, f func()) {
	clientLock.Lock()
	orig := client
	client = papiCli

	origHapi := hapiClient
	hapiClient = hapiCli

	defer func() {
		client = orig
		hapiClient = origHapi
		clientLock.Unlock()
	}()

	f()
}

// suppressLogging prevents logging output during the given func unless TEST_LOGGING env var is not empty. Use this
// to keep log messages from polluting test output. Not thread-safe.
func suppressLogging(t *testing.T, f func()) {
	t.Helper()

	if os.Getenv("TEST_LOGGING") == "" {
		orig := hclog.SetDefault(hclog.NewNullLogger())
		defer func() { hclog.SetDefault(orig) }()
		t.Log("Logging is suppressed. Set TEST_LOGGING=1 in env to see logged messages during test")
	}

	f()
}

// Wrapper to intercept the papi.Mock's call of t.FailNow(). The Terraform test driver runs the provider code on
// goroutines other than the one created for the test. When t.FailNow() is called from any other goroutine, it causes
// the test to hang because the TF test driver is still waiting to serve requests. Mockery's failure message neglects to
// inform the user which test had failed. Use this struct to wrap a *testing.T when you call mock.Test(T{t}) and the
// mock's failure will print the failling test's name. Such failures are usually caused by the provider invoking an
// unexpected call on the mock.
//
// NB: You should only need to use this where your test uses the Terraform test driver
type T struct{ *testing.T }

// Overrides testing.T.FailNow() so when a test mock fails an assertion, we see which test had failed before it hangs
func (t T) FailNow() {
	t.T.Fatalf("FAIL: %s", t.T.Name())
}

type (
	TestSubprovider struct {
		resources   []func() resource.Resource
		datasources []func() datasource.DataSource
		client      papi.PAPI
	}

	clientSetter interface {
		setClient(papi.PAPI)
	}

	testSubproviderOption func(*TestSubprovider)
)

func withMockClient(mock papi.PAPI) testSubproviderOption {
	return func(ts *TestSubprovider) {
		ts.client = mock
	}
}

func newTestSubprovider(opts ...testSubproviderOption) *TestSubprovider {
	s := NewFrameworkSubprovider()

	ts := &TestSubprovider{
		resources:   s.Resources(),
		datasources: s.DataSources(),
	}

	for _, opt := range opts {
		opt(ts)
	}

	return ts
}

// Resources returns terraform resources for property
func (ts *TestSubprovider) Resources() []func() resource.Resource {
	for i, fn := range ts.resources {
		// decorate
		fn := fn
		ts.resources[i] = func() resource.Resource {
			res := fn()
			if v, ok := res.(clientSetter); ok {
				v.setClient(ts.client)
			}
			return res
		}
	}
	return ts.resources
}

// DataSources returns terraform data sources for property
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
