package property

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataPropertyRules(t *testing.T) {
	tests := map[string]struct {
		givenTF  string
		expectJS string
	}{
		"siteshield": {
			givenTF:  "siteshield.tf",
			expectJS: "siteshield.json",
		},
		"cpcode": {
			givenTF:  "cpCode.tf",
			expectJS: "cpCode.json",
		},
		"criteria": {
			givenTF:  "criteria.tf",
			expectJS: "criteria.json",
		},
		"is secure false": {
			givenTF:  "isSecureFalse.tf",
			expectJS: "isSecureFalse.json",
		},
		"is secure true": {
			givenTF:  "isSecureTrue.tf",
			expectJS: "isSecureTrue.json",
		},
		"variables": {
			givenTF:  "variables.tf",
			expectJS: "variables.json",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			expectJS := compactJSON(loadFixtureBytes(fmt.Sprintf("testdata/TestDataPropertyRules/%s", test.expectJS)))
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Destroy: false,
						Config:  loadFixtureString(fmt.Sprintf("testdata/TestDataPropertyRules/%s", test.givenTF)),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", expectJS),
							resource.TestCheckResourceAttrSet("data.akamai_property_rules.rules", "json"),
						),
					},
				},
			})
		})
	}
}
