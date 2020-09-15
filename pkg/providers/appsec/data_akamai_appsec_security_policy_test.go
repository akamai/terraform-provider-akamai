package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiSecurityPolicy_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_security_policy.appsecsecuritypolicy"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSecurityPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiSecurityPolicyConfig() string {
	return `
provider "akamai" {
  appsec_section = "default"
}


data "akamai_appsec_configuration_version" "appsecconfigurationversion" {
    name = "Akamai Tools"
   }

output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}

output "configsedgeconfiglist" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_list
}

output "configsedgeconfigversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.version
}
data "akamai_appsec_security_policy" "appsecsecuritypolicy" {
  name = "akamaitools" 
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version =  data.akamai_appsec_configuration.appsecconfigedge.version
}

output "securitypolicy" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.policy_id
}

output "securitypolicies" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.policy_list
}

`
}

func testAccCheckAkamaiSecurityPolicyDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai SecurityPolicy] SecurityPolicy Destroy skipped ")
	return nil
}
