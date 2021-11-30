provider "akamai" {
  edgerc = "~/.edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_asmap" "tfexample_as_1" {
  domain = local.gtmTestDomain
  name   = "tfexample_as_1"
  default_datacenter {
    datacenter_id = 5400
    nickname      = "default datacenter"
  }
  assignment {
    datacenter_id = 3132
    nickname      = "tfexample_dc_2"
    as_numbers    = [12223, 16701, 17333]
  }
  assignment {
    datacenter_id = 3133
    nickname      = "tfexample_dc_3"
    as_numbers    = [12228, 16704, 17336]
  }
  wait_on_complete = true
}

