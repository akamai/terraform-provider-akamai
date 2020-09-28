package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiRatePolicies_basic(t *testing.T) {
	dataSourceName := "data.appsec_akamai_rate_policies.akamairatepolicies"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiRatePoliciesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiRatePoliciesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiRatePoliciesConfig() string {
	return `
provider "akamai" {
  appsec_section = "default"
}
data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}


data "akamai_appsec_rate_policies" "appsecreatepolicies" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version_number = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}


`
}

func testAccCheckAkamaiRatePoliciesDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai RatePolicies] RatePolicies Destroy skipped ")
	return nil
}
