package cps

import (
	"fmt"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider

var testProvider *schema.Provider

func init() {
	testProvider = akamai.Provider(Subprovider())()
	testAccProviders = map[string]*schema.Provider{
		"akamai": testProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client cps.CPS, f func()) {
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
