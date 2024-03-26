provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
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
    datacenter_id = 3131
    nickname      = "tfexample_dc_1"
    // Optional
    countries = ["WE", "US", "FR", "GB"]
  }
  assignment {
    datacenter_id = 3132
    nickname      = "tfexample_dc_2"
    // Optional
    countries = ["GB", "AU"]
  }
  assignment {
    datacenter_id = 3133
    nickname      = "tfexample_dc_3"
    // Optional
    countries = ["GB", "BG", "CN", "MC", "TR"]
  }
  wait_on_complete = true
}

