package edgeworkers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_edgeworkers_property_rules.test", "json",
								loadFixtureString(test.expectedJSONPath)),
							resource.TestCheckResourceAttr(
								"data.akamai_edgeworkers_property_rules.test", "id", "123"),
						),
					},
				},
			})
		})
	}
}
