package property

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/go-hclog"
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

func testAccPreCheck(t *testing.T) {
	TODO(t, "Check not implemented")
}

func getTestProvider() *schema.Provider {
	return testProvider
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client papi.PAPI, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	defer func() {
		inst.client = orig
		clientLock.Unlock()
	}()

	f()
}

// TODO marks a test as being in a "pending" state and logs a message telling the user why. Such tests are expected to
// fail for the time being and may exist for the sake of unfinished/future features or to document known buggy cases
// that won't be fixed right away. The failure of a pending test is not considered an error and the test will therefore
// be skipped unless the TEST_TODO environment variable is set to a non-empty value.
func TODO(t *testing.T, message string) {
	t.Helper()
	t.Log(fmt.Sprintf("TODO: %s", message))

	if os.Getenv("TEST_TODO") == "" {
		t.Skip("TODO: Set TEST_TODO=1 in env to run this test")
	}
}

func setEnv(home string, env map[string]string) {
	os.Clearenv()
	os.Setenv("HOME", home)
	if len(env) > 0 {
		for key, val := range env {
			os.Setenv(key, val)
		}
	}
}

func restoreEnv(env []string) {
	os.Clearenv()
	for _, value := range env {
		envVar := strings.Split(value, "=")
		os.Setenv(envVar[0], envVar[1])
	}
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

// compactJSON converts a JSON-encoded byte slice to a compact form (so our JSON fixtures can be readable)
func compactJSON(encoded []byte) string {
	buf := bytes.Buffer{}
	if err := json.Compact(&buf, encoded); err != nil {
		panic(fmt.Sprintf("%s: %s", err, string(encoded)))
	}

	return buf.String()
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
