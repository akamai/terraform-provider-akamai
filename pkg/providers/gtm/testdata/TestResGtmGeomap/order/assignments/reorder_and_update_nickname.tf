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
    datacenter_id = 3133
    nickname      = "tfexample_dc_30"
    // Optional
    countries = ["GB", "BG", "CN", "MC", "TR"]
  }
  assignment {
    datacenter_id = 3131
    nickname      = "tfexample_dc_1"
    // Optional
    countries = ["GB", "PL", "US", "FR"]
  }
  assignment {
    datacenter_id = 3132
    nickname      = "tfexample_dc_2"
    // Optional
    countries = ["GB", "AU"]
  }
  wait_on_complete = true
}

