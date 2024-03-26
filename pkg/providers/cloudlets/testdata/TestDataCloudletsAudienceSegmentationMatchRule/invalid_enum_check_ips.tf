provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      check_ips      = "incorrect_check_ips"
      object_match_value {
        type  = "simple"
        value = ["a"]
      }
    }
    forward_settings {
    }
  }
}