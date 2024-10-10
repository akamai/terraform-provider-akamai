provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_datacenter" "tfexample_dc_1" {
  domain           = local.gtmTestDomain
  nickname         = "tfexample_dc_1"
  wait_on_complete = false
  city             = "Snæfellsjökull"
  country          = "IS"
  latitude         = "64.808"
  longitude        = "-23.776"
  default_load_object {
    load_object      = "/test"
    load_object_port = 80
    load_servers     = ["1.2.3.4", "1.2.3.9"]
  }
  continent = "EU"
}
