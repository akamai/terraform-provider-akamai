provider "akamai" {
  edgerc = "~/.edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
  contract = "9-CONTRACT"
  group    = "12345"
}

/*
resource "akamai_gtm_domain" "testdomain" {
        name = "gtm_terra_testdomain.akadns.net"
        type = "weighted"
        contract = "1-2ABCDEF"
        comment =  "Test"
        group     = "123ABC" 
        load_imbalance_percentage = 10
}
*/

resource "akamai_gtm_datacenter" "tfexample_dc_1" {
  domain = local.gtmTestDomain
  nickname         = "tfexample_dc_1"
  wait_on_complete = false
  default_load_object {
    load_object      = "/test"
    load_object_port = 80
    load_servers     = ["1.2.3.4", "1.2.3.9"]
  }
}
