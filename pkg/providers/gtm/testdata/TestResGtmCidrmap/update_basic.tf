provider "akamai" {
  edgerc = "../../test/edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_cidrmap" "tfexample_cidrmap_1" {
  domain = local.gtmTestDomain
  name   = "tfexample_cidrmap_1"
  default_datacenter {
    datacenter_id = 5400
    nickname      = "default datacenter"
  }
  assignment {
    datacenter_id = 3132
    nickname      = "tfexample_dc_2"
    // Optional
    blocks = ["1.2.3.9/16"]
  }
  wait_on_complete = true
}

