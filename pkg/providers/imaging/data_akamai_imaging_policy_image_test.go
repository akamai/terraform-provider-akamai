package imaging

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataPolicyImage(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"empty policy": {
			configPath:       "testdata/TestDataPolicyImage/empty_policy/policy.tf",
			expectedJSONPath: "testdata/TestDataPolicyImage/empty_policy/expected.json",
		},
		"regular policy with 1 transformation": {
			configPath:       "testdata/TestDataPolicyImage/regular_policy/policy.tf",
			expectedJSONPath: "testdata/TestDataPolicyImage/regular_policy/expected.json",
		},
		"regular policy with multiple nested transformations": {
			configPath:       "testdata/TestDataPolicyImage/policy_with_nested_transformations/policy.tf",
			expectedJSONPath: "testdata/TestDataPolicyImage/policy_with_nested_transformations/expected.json",
		},
		"policy with empty breakpoints": {
			configPath:       "testdata/TestDataPolicyImage/empty_breakpoints/policy.tf",
			expectedJSONPath: "testdata/TestDataPolicyImage/empty_breakpoints/expected.json",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if name == "policy with empty breakpoints" {
				t.Skip("It should be restored once DXE-941 is fixed")
			}
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_imaging_policy_image.policy", "json",
								loadFixtureString(test.expectedJSONPath)),
						),
					},
				},
			})
		})
	}
}
