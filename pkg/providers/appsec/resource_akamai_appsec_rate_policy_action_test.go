package appsec

import (
	"strconv"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiRatePolicyAction_basic(t *testing.T) {

	dataSourceName := "akamai_appsec_rate_policy_action.appsecratepolicyaction"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAkamaiRatePolicyActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiRatePolicyActionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiRatePolicyActionExists,
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiRatePolicyActionConfig() string {
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


resource "akamai_appsec_rate_policy_action" "appsecreatepolicyaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "f1rQ_106946"
    ipv4_action = "alert"
    ipv6_action = "none"
}

output "ratepolicyaction" {
  value = akamai_appsec_rate_policy_action.appsecreatepolicyaction.rate_policy_id
}


`
}

func testAccCheckAkamaiRatePolicyActionExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_rate_policy_action" {
			continue
		}
		//rname := rs.Primary.ID
		configid, _ := strconv.Atoi(rs.Primary.Attributes["config_id"])
		version, _ := strconv.Atoi(rs.Primary.Attributes["version"])
		policyid := rs.Primary.Attributes["policy_id"]
		//rate_policy_id, _ := strconv.Atoi(rs.Primary.Attributes["rate_policy_id"])

		ccresp := appsec.NewRatePolicyActionResponse()
		err := ccresp.GetRatePolicyAction(configid, version, policyid, "TEST")

		if err != nil {
			return err
		}
	}

	return nil
}
