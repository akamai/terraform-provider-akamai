provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      object_match_value {
        type  = "incorrect_type"
        value = ["a"]
      }
    }
    forward_settings {
    }
  }
}