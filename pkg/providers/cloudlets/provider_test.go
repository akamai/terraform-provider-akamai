package cloudlets

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets/v3"
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
			upgradedPluginProvider, err := tf5to6server.UpgradeServer(
				context.Background(),
				testAccPluginProvider.GRPCProvider,
			)
			if err != nil {
				return nil, err
			}
			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return upgradedPluginProvider
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

// useClientV2AndV3 swaps out both client (v2) and client v3 on the global instances for the duration of the given func. To be used in by tests for data sources and resources that use both V2 & V3 cloudlets
func useClientV2AndV3(cloudletsV2Client cloudlets.Cloudlets, cloudletsV3Client v3.Cloudlets, f func()) {
	clientLock.Lock()
	origV2 := client
	client = cloudletsV2Client
	origV3 := v3Client
	v3Client = cloudletsV3Client

	defer func() {
		client = origV2
		v3Client = origV3
		clientLock.Unlock()
	}()

	f()
}
