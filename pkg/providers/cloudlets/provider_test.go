package cloudlets

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProviders         map[string]func() (tfprotov5.ProviderServer, error)
	testAccPluginProvider    *schema.Provider
	testAccFrameworkProvider provider.Provider
)

func TestMain(m *testing.M) {
	testAccPluginProvider = akamai.NewPluginProvider(NewPluginSubprovider())()
	testAccFrameworkProvider = akamai.NewFrameworkProvider(NewFrameworkSubprovider())()

	testAccProviders = map[string]func() (tfprotov5.ProviderServer, error){
		"akamai": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				testAccPluginProvider.GRPCProvider,
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
func useClient(cloudletsClient cloudlets.Cloudlets, f func()) {
	clientLock.Lock()
	orig := client
	client = cloudletsClient

	defer func() {
		client = orig
		clientLock.Unlock()
	}()

	f()
}

// useClientV3 swaps out the client v3 on the global instance for the duration of the given func
func useClientV3(cloudletsV3Client v3.Cloudlets, f func()) {
	clientLock.Lock()
	orig := v3Client
	v3Client = cloudletsV3Client

	defer func() {
		v3Client = orig
		clientLock.Unlock()
	}()

	f()
}
