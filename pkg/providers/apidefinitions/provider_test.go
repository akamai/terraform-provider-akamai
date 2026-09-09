package apidefinitions

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

var providerConfig = `provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}
`
