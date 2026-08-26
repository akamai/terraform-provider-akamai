provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_property" "tfexample_prop_1" {
  domain                      = local.gtmTestDomain
  name                        = "tfexample_prop_1"
  type                        = "failover"
  ipv6                        = true
  score_aggregation_type      = "median"
  stickiness_bonus_percentage = 10
  stickiness_bonus_constant   = 10
  health_threshold            = 123
  use_computed_targets        = true
  backup_ip                   = "test ip"
  balance_by_download_score   = true
  unreachable_threshold       = 1234
  min_live_fraction           = 1
  health_multiplier           = 5
  dynamic_ttl                 = 300
  max_unreachable_penalty     = 123
  map_name                    = "test map"
  handout_limit               = 5
  handout_mode                = "normal"
  failover_delay              = 5
  backup_cname                = "test cname"
  failback_delay              = 5
  load_imbalance_percentage   = 10
  health_max                  = 123
  ghost_demand_reporting      = false
  cname                       = "test cName"
  comments                    = "test comment"

  traffic_target {
    datacenter_id = 3131
    enabled       = true
    weight        = 200
    servers       = ["1.2.3.9"]
    handout_cname = "test"
    precedence    = 10
  }

  static_rr_set {
    type  = "MX"
    ttl   = 300
    rdata = ["100 test_e"]
  }

  liveness_test {
    name                             = "lt5"
    test_interval                    = 40
    test_object_protocol             = "HTTP"
    test_timeout                     = 30
    test_object                      = "/junk"
    test_object_port                 = 1
    disable_nonstandard_port_warning = false
    http_header {
      name  = "test_name"
      value = "test_value"
    }
  }

  state_change_notification_webhook {
    url    = "https://example.com/gtm-webhook"
    format = "json-compact"
  }

  wait_on_complete = false
}
