package botman

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
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
	getLatestConfigVersion = func(_ context.Context, _ int, _ interface{}) (int, error) {
		return 15, nil
	}
	getModifiableConfigVersion = func(_ context.Context, _ int, _ string, _ interface{}) (int, error) {
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
