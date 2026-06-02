package clientlists

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// clientListTestCase is a shared test-case struct used by data-source tests.
type clientListTestCase struct {
	init  func(*clientlists.Mock)
	steps []resource.TestStep
}

// runClientListTestCases executes a map of clientListTestCase subtests.
func runClientListTestCases(t *testing.T, tests map[string]clientListTestCase) {
	t.Helper()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &clientlists.Mock{}
			if tc.init != nil {
				tc.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

// runResourceTest is a convenience wrapper for resource subtests.
func runResourceTest(t *testing.T, client *clientlists.Mock, steps []resource.TestStep) {
	t.Helper()
	useClient(client, func() {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps:                    steps,
		})
	})
	client.AssertExpectations(t)
}
