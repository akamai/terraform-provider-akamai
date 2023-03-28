package botman

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func TestMain(m *testing.M) {
	testAccProvider = akamai.Provider(Subprovider())()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"akamai": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
	if err := akamai.TFTestSetup(); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := akamai.TFTestTeardown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
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
