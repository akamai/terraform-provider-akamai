provider "akamai" {
  edgerc = "~/.edgerc"
}

locals {  
  contract = "9-CONTRACT"
  group    = "12345"
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

/*
resource "akamai_gtm_domain" "testdomain" {
        name = "gtm_terra_testdomain.akadns.net"
        type = "weighted"
        contract = "1-2ABCDEF"
        comment =  "Test"
        group     = "123ABC" 
        load_imbalance_percentage = 10
}
*/

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
    name          = ""
    handout_cname = "test"
  }

  liveness_test {
    name                             = "lt5"
    test_interval                    = 40
    test_object_protocol             = "HTTP"
    test_timeout                     = 30
    answers_required                 = false
    disable_nonstandard_port_warning = false
    error_penalty                    = 0
    http_error3xx                    = false
    http_error4xx                    = false
    http_error5xx                    = false
    disabled                         = false
    http_header {
      name  = "test_name"
      value = "test_value"
    }
    peer_certificate_verification = false
    recursion_requested           = false
    request_string                = ""
    resource_type                 = ""
    response_string               = ""
    ssl_client_certificate        = ""
    ssl_client_private_key        = ""
    test_object                   = "/junk"
    test_object_password          = ""
    test_object_port              = 1
    test_object_username          = ""
    timeout_penalty               = 0
  }
  liveness_test {
    name                 = "lt2"
    test_interval        = 30
    test_object_protocol = "HTTP"
    test_timeout         = 20
    test_object          = "/junk"
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

