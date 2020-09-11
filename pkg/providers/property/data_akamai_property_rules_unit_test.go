package property

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAkamaiDataPropertyRules_siteshield(t *testing.T) {
	expectJS := compactJSON(loadFixtureBytes("testdata/TestAkamaiDataPropertyRules_siteshield.json"))

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  loadFixtureString("testdata/TestAkamaiDataPropertyRules_siteshield.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", expectJS),
					resource.TestCheckResourceAttrSet("data.akamai_property_rules.rules", "json"),
				),
			},
		},
	})
}

func TestAkamaiDataPropertyRules_cpCode(t *testing.T) {
	expectJS := compactJSON(loadFixtureBytes("testdata/TestAkamaiDataPropertyRules_cpCode.json"))

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  loadFixtureString("testdata/TestAkamaiDataPropertyRules_cpCode.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", expectJS),
					resource.TestCheckResourceAttrSet("data.akamai_property_rules.rules", "json"),
				),
			},
		},
	})
}

func TestAkamaiDataPropertyRules_isSecureTrue(t *testing.T) {
	expectJS := compactJSON(loadFixtureBytes("testdata/TestAkamaiDataPropertyRules_isSecureTrue.json"))

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  loadFixtureString("testdata/TestAkamaiDataPropertyRules_isSecureTrue.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", expectJS),
					resource.TestCheckResourceAttrSet("data.akamai_property_rules.rules", "json"),
				),
			},
		},
	})
}

func TestAkamaiDataPropertyRules_isSecureFalse(t *testing.T) {
	expectJS := compactJSON(loadFixtureBytes("testdata/TestAkamaiDataPropertyRules_isSecureFalse.json"))

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  loadFixtureString("testdata/TestAkamaiDataPropertyRules_isSecureFalse.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", expectJS),
					resource.TestCheckResourceAttrSet("data.akamai_property_rules.rules", "json"),
				),
			},
		},
	})
}
