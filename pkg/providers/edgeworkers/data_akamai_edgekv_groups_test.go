package edgeworkers

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEdgeKVGroups(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		client.EdgeWorkers.On("ListGroupsWithinNamespace", testutils.MockContext, edgeworkers.ListGroupsWithinNamespaceRequest{
			Network:     "staging",
			NamespaceID: "test_namespace"}).
			Return([]string{"TestImportGroup", "TestGroup1", "TestGroup2", "TestGroup3", "TestGroup4"}, nil).Times(3)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVNamespaceGroups/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_edgekv_groups.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_edgekv_groups.test", "id", "test_namespace:staging"),
						resource.TestCheckResourceAttr("data.akamai_edgekv_groups.test", "groups.#", "5"),
						resource.TestCheckResourceAttr("data.akamai_edgekv_groups.test", "groups.0", "TestImportGroup"),
						resource.TestCheckResourceAttr("data.akamai_edgekv_groups.test", "groups.1", "TestGroup1"),
						resource.TestCheckResourceAttr("data.akamai_edgekv_groups.test", "groups.2", "TestGroup2"),
						resource.TestCheckResourceAttr("data.akamai_edgekv_groups.test", "groups.3", "TestGroup3"),
						resource.TestCheckResourceAttr("data.akamai_edgekv_groups.test", "groups.4", "TestGroup4"),
					),
				},
			},
		})

		client.EdgeWorkers.AssertExpectations(t)
	})

	t.Run("missed required `namespace_name` field", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVNamespaceGroups/missed_namespace_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "namespace_name" is required, but no definition was found.`),
				},
			},
		})

		client.EdgeWorkers.AssertExpectations(t)
	})

	t.Run("missed required `network` field", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVNamespaceGroups/missed_network.tf"),
					ExpectError: regexp.MustCompile(`The argument "network" is required, but no definition was found.`),
				},
			},
		})

		client.EdgeWorkers.AssertExpectations(t)
	})

	t.Run("incorrect `network` field", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVNamespaceGroups/incorrect_network.tf"),
					ExpectError: regexp.MustCompile(`expected network to be one of \["staging" "production"], got incorrect_network`),
				},
			},
		})

		client.EdgeWorkers.AssertExpectations(t)
	})
}
