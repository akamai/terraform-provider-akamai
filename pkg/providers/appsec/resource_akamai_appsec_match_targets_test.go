package appsec

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiMatchTargets_basic(t *testing.T) {
	dataSourceName := "akamai_appsec_match_targets.appsecmatchtargets"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiMatchTargetsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiMatchTargetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func TestAccAkamaiMatchTargets_update(t *testing.T) {
	dataSourceName := "akamai_appsec_match_targets.appsecmatchtargets"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiMatchTargetsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiMatchTargetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					testAccCheckAkamaiMatchTargetsExists,
					//resource.TestCheckResourceAttr("akamai_appsec_match_targets.appsecmatchtargets", "load_imbalance_percentage", "50"),
				),
			},
			{
				Config: testAccAkamaiMatchTargetsUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					testAccCheckAkamaiMatchTargetsExists,
					//resource.TestCheckResourceAttr("akamai_appsec_match_targets.appsecmatchtargets", "load_imbalance_percentage", "50"),
				),
			},
		},
	})
}

func testAccAkamaiMatchTargetsConfig() string {
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


resource "akamai_appsec_match_targets" "appsecmatchtargets" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    type =  "website"
    sequence =  1
    is_negative_path_match =  false
    is_negative_file_extension_match =  true
    default_file = "BASE_MATCH"
    hostnames =  ["example.com","www.example.net","m.example.com"]
    //file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
    //file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
    security_policy = "f1rQ_106946"
 
    bypass_network_lists = ["888518_ACDDCKERS","1304427_AAXXBBLIST"]
    
}

`
}

func testAccAkamaiMatchTargetsUpdateConfig() string {
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


resource "akamai_appsec_match_targets" "appsecmatchtargets" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    type =  "website"
    sequence =  1
    is_negative_path_match =  false
    is_negative_file_extension_match =  true
    default_file = "BASE_MATCH"
    hostnames =  ["example.com","www.example.net","p.example.com"]
    //file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
    //file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
    security_policy = "f1rQ_106946"
 
    bypass_network_lists = ["888518_ACDDCKERS","1304427_AAXXBBLIST"]
    
}

`
}

func testCheckDeleteMatchTargetsResource(s *terraform.State, rscName string) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_match_targets" {
			continue
		}

		ccresp := appsec.NewMatchTargetsResponse()

		ccresp.ConfigID, _ = strconv.Atoi(rs.Primary.Attributes["config_id"])
		ccresp.ConfigVersion, _ = strconv.Atoi(rs.Primary.Attributes["version"])
		ccresp.TargetID, _ = strconv.Atoi(rs.Primary.ID)

		err := ccresp.GetMatchTargets("TEST")

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Test] Deleting test resource [%v]", rscName)

		err = ccresp.DeleteMatchTargets("TEST")
		if err != nil {
			return fmt.Errorf("resource was not deleted %s. Error: %s", rscName, err.Error())
		}
	}

	return nil
}

func testAccCheckAkamaiMatchTargetsDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_match_targets" {
			continue
		}

		rscName := "akamai_appsec_match_targets"
		if err := testCheckDeleteMatchTargetsResource(s, rscName); err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckAkamaiMatchTargetsExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_match_targets" {
			continue
		}

		//rname := rs.Primary.ID

		ccresp := appsec.NewMatchTargetsResponse()
		ccresp.ConfigID, _ = strconv.Atoi(rs.Primary.Attributes["config_id"])
		ccresp.ConfigVersion, _ = strconv.Atoi(rs.Primary.Attributes["version"])
		ccresp.TargetID, _ = strconv.Atoi(rs.Primary.ID)
		err := ccresp.GetMatchTargets("TEST")

		if err != nil {
			return err
		}
	}

	return nil
}
