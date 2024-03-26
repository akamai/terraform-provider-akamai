provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
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
    as_numbers    = [16703, 12229, 17335]
  }
  assignment {
    datacenter_id = 3133
    nickname      = "tfexample_dc_3"
    as_numbers    = [2222, 5555, 1111, 3333, 4444]
  }
  assignment {
    datacenter_id = 3131
    nickname      = "tfexample_dc_1"
    as_numbers    = [17334, 16702, 12222]
  }
  wait_on_complete = true
}

