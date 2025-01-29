package imaging

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		"policy with composite_post_policy breakpoints": {
			configPath:       "testdata/TestDataPolicyImage/composite_post_policy/policy.tf",
			expectedJSONPath: "testdata/TestDataPolicyImage/composite_post_policy/expected.json",
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
								"data.akamai_imaging_policy_image.policy", "json",
								testutils.LoadFixtureString(t, test.expectedJSONPath)),
						),
					},
				},
			})
		})
	}
}
