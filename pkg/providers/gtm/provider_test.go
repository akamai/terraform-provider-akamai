package gtm

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvidersProtoV5 map[string]func() (tfprotov5.ProviderServer, error)
var testAccProviders map[string]func() (*schema.Provider, error)
var testAccFrameworkProvider provider.Provider
var testAccProvider *schema.Provider

func TestMain(m *testing.M) {
	testAccProvider = akamai.NewPluginProvider(NewSubprovider())()
	testAccFrameworkProvider = akamai.NewFrameworkProvider(NewFrameworkSubprovider())()

	testAccProviders = map[string]func() (*schema.Provider, error){
		"akamai": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
	testAccProvidersProtoV5 = map[string]func() (tfprotov5.ProviderServer, error){
		"akamai": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				providerserver.NewProtocol5(
					testAccFrameworkProvider,
				),
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
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
func useClient(client gtm.GTM, f func()) {
	clientLock.Lock()
	orig := inst.client
	origFrameworkClient := frameworkInst.client
	inst.client = client
	frameworkInst.client = client
	defer func() {
		inst.client = orig
		frameworkInst.client = origFrameworkClient
		clientLock.Unlock()
	}()

	f()
}
