provider "akamai" {
  edgerc = "~/.edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_geomap" "tfexample_geomap_1" {
  domain = local.gtmTestDomain
  name   = "tfexample_geomap_1"
  default_datacenter {
    datacenter_id = 5400
    nickname      = "default datacenter"
  }
  assignment {
    datacenter_id = 3132
    nickname      = "tfexample_dc_2"
    // Optional
    countries = ["US"]
  }
  wait_on_complete = true
}

