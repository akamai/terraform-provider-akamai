package appsec

import (
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiSecurityPolicyClone_basic(t *testing.T) {
	dataSourceName := "akamai_appsec_security_policy_clone.appsecsecuritypolicyclone"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAkamaiSecurityPolicyCloneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSecurityPolicyCloneConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiSecurityPolicyCloneConfig() string {
	return `
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}


resource "akamai_appsec_security_policy_clone" "appsecsecuritypolicyclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    
    create_from_security_policy = "LNPD_76189"
    policy_name = "Cloned Test for Launchpad 22"
    policy_prefix = "LN" 
   }

output "secpolicyclone" {
  value = akamai_appsec_security_policy_clone.appsecsecuritypolicyclone.policy_id
}

`
}

func testAccCheckAkamaiSecurityPolicyCloneExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_security_policy_clone" {
			continue
		}
		//rname := rs.Primary.ID
		ccresp := appsec.NewSecurityPolicyCloneResponse()
		_, err := ccresp.GetSecurityPolicyClone("TEST")

		if err != nil {
			return err
		}
	}
	return nil
}
