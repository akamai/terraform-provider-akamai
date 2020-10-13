provider "akamai" {
  edgerc = "~/.edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_datacenter" "tfexample_dc_1" {
  domain = local.gtmTestDomain
  nickname         = "tfexample_dc_1"
  wait_on_complete = true 
  default_load_object {
    load_object      = "/test"
    load_object_port = 80
    load_servers     = ["1.2.3.5", "1.2.3.6"]
  }
  country = "US"
}
