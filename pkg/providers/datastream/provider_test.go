package datastream

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)

var testAccProvider *schema.Provider

func TestMain(m *testing.M) {
	testAccProvider = akamai.NewPluginProvider(Subprovider())()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"akamai": func() (*schema.Provider, error) {
			return testAccProvider, nil
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

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client datastream.DS, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	defer func() {
		inst.client = orig
		clientLock.Unlock()
	}()

	f()
}

// loadFixtureBytes returns the entire contents of the given file as a byte slice
func loadFixtureBytes(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// loadFixtureString returns the entire contents of the given file as a string
func loadFixtureString(format string, args ...interface{}) string {
	return string(loadFixtureBytes(fmt.Sprintf(format, args...)))
}
