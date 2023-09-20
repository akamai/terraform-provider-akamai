package botman

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

var testAccProviders map[string]func() (tfprotov6.ProviderServer, error)

func TestMain(m *testing.M) {
	testAccSDKProvider := akamai.NewSDKProvider(NewSubprovider())()
	testAccProviders = map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			sdkProviderV6, err := tf5to6server.UpgradeServer(
				context.Background(),
				testAccSDKProvider.GRPCProvider,
			)
			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return sdkProviderV6
				},
			}

			muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)
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
func useClient(client *botman.Mock, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client
	origGetLatestConfigVersion := getLatestConfigVersion
	origGetModifiableConfigVersion := getModifiableConfigVersion
	getLatestConfigVersion = func(ctx context.Context, configID int, m interface{}) (int, error) {
		return 15, nil
	}
	getModifiableConfigVersion = func(ctx context.Context, configID int, resource string, m interface{}) (int, error) {
		return 15, nil
	}
	defer func() {
		inst.client = orig
		getLatestConfigVersion = origGetLatestConfigVersion
		getModifiableConfigVersion = origGetModifiableConfigVersion
		clientLock.Unlock()
	}()
	f()
}

func compactJSON(message string) string {
	var dst bytes.Buffer
	err := json.Compact(&dst, []byte(message))
	if err != nil {
		panic(err)
	}
	return dst.String()
}
