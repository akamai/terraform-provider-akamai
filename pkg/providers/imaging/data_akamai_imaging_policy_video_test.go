package imaging

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataPolicyVideo(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"empty policy": {
			configPath:       "testdata/TestDataPolicyVideo/empty_policy/policy.tf",
			expectedJSONPath: "testdata/TestDataPolicyVideo/empty_policy/expected.json",
		},
		"regular policy": {
			configPath:       "testdata/TestDataPolicyVideo/regular_policy/policy.tf",
			expectedJSONPath: "testdata/TestDataPolicyVideo/regular_policy/expected.json",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_imaging_policy_video.policy", "json",
								testutils.LoadFixtureString(t, test.expectedJSONPath)),
						),
					},
				},
			})
		})
	}
}
