package imaging

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
)

func TestMain(m *testing.M) {
	PolicyDepth = 4
	testutils.TestRunner(m)
}
