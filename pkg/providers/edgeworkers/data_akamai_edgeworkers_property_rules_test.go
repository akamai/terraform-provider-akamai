package edgeworkers

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataEdgeworkersPropertyRules(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"with provided edgeworker ID": {
			configPath:       "testdata/TestDataEdgeWorkersPropertyRules/basic.tf",
			expectedJSONPath: "testdata/TestDataEdgeWorkersPropertyRules/rules/basic.json",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_edgeworkers_property_rules.test", "json",
								testutils.LoadFixtureString(t, test.expectedJSONPath)),
							resource.TestCheckResourceAttr(
								"data.akamai_edgeworkers_property_rules.test", "id", "123"),
						),
					},
				},
			})
		})
	}
}
