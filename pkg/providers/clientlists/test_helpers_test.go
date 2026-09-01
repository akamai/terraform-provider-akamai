package clientlists

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// clientListTestCase is a shared test-case struct used by both resource and
// data-source tests.
type clientListTestCase struct {
	init  func(*clientlists.Mock)
	steps []resource.TestStep
}

// runClientListFrameworkTestCases executes a map of clientListTestCase subtests against
// the Framework provider factory (used by data sources).
func runClientListFrameworkTestCases(t *testing.T, tests map[string]clientListTestCase) {
	t.Helper()
	runTestCases(t, tests, true)
}

// runClientListSDKTestCases executes a map of clientListTestCase subtests against the
// SDKv2 provider factory (used by resources).
func runClientListSDKTestCases(t *testing.T, tests map[string]clientListTestCase) {
	t.Helper()
	runTestCases(t, tests, false)
}

// runTestCases runs each clientListTestCase as a parallel subtest. When
// useFrameworkProvider is true the Framework provider factory is used,
// otherwise the SDKv2 provider factory is used.
func runTestCases(t *testing.T, tests map[string]clientListTestCase, useFrameworkProvider bool) {
	t.Helper()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.ClientLists)
			}

			providerFactories := testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig()))
			if useFrameworkProvider {
				providerFactories = testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig()))
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: providerFactories,
				Steps:                    tc.steps,
			})

			client.ClientLists.AssertExpectations(t)
		})
	}
}
