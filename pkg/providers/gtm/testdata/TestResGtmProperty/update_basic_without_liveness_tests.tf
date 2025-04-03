provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_property" "tfexample_prop_1" {
  domain                 = local.gtmTestDomain
  name                   = "tfexample_prop_1"
  type                   = "weighted-round-robin"
  score_aggregation_type = "median"
  handout_limit          = 5
  handout_mode           = "normal"
  traffic_target {
    datacenter_id = 3131
    enabled       = true
    weight        = 200
    servers       = ["1.2.3.9"]
    handout_cname = "test"
  }
  static_rr_set {
    type  = "MX"
    ttl   = 300
    rdata = ["100 test_e"]
  }
  failover_delay   = 0
  failback_delay   = 0
  wait_on_complete = false
}