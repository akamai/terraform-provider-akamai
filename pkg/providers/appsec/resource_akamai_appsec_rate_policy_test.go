package appsec

import (
	"fmt"
	"log"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiRatePolicy_basic(t *testing.T) {
	dataSourceName := "appsec_appsec_rate_policy.appsecratepolicy"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiRatePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiRatePolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func TestAccAkamaiRatePolicy_update(t *testing.T) {
	dataSourceName := "appsec_appsec_rate_policy.appsecratepolicy"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiRatePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiRatePolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					testAccCheckAkamaiRatePolicyExists,
					//resource.TestCheckResourceAttr("appsec_appsec_rate_policy.appsecratepolicy", "load_imbalance_percentage", "50"),
				),
			},
			{
				Config: testAccAkamaiRatePolicyUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					testAccCheckAkamaiRatePolicyExists,
					//resource.TestCheckResourceAttr("appsec_appsec_rate_policy.appsecratepolicy", "load_imbalance_percentage", "50"),
				),
			},
		},
	})
}

func testAccAkamaiRatePolicyConfig() string {
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


resource "akamai_appsec_rate_policy" "appsecratepolicy" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}


`
}

func testAccAkamaiRatePolicyUpdateConfig() string {
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


resource "akamai_appsec_rate_policy" "appsecratepolicy" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}


`
}

func testCheckDeleteRatePolicyResource(s *terraform.State, rscName string) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_rate_policy" {
			continue
		}

		ccresp := appsec.NewRatePolicyResponse()

		err := ccresp.GetRatePolicy("TEST")

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Test] Deleting test resource [%v]", rscName)

		err = ccresp.DeleteRatePolicy("TEST")
		if err != nil {
			return fmt.Errorf("resource was not deleted %s. Error: %s", rscName, err.Error())
		}
	}

	return nil
}

func testAccCheckAkamaiRatePolicyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_rate_policy" {
			continue
		}

		rscName := "akamai_appsec_rate_policy"
		if err := testCheckDeleteRatePolicyResource(s, rscName); err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckAkamaiRatePolicyExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_rate_policy" {
			continue
		}

		//rname := rs.Primary.ID

		ccresp := appsec.NewRatePolicyResponse()

		err := ccresp.GetRatePolicy("TEST")

		if err != nil {
			return err
		}
	}

	return nil
}
