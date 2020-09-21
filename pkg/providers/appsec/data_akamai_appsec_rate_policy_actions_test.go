package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiRatePolicyActions_basic(t *testing.T) {
	dataSourceName := "data.appsec_akamai_rate_policy_actions.akamairatepolicyactions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiRatePolicyActionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiRatePolicyActionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiRatePolicyActionsConfig() string {
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


data "akamai_appsec_rate_policy_actions" "appsecreatepolicysaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}


`
}

func testAccCheckAkamaiRatePolicyActionsDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai RatePolicyActions] RatePolicyActions Destroy skipped ")
	return nil
}
