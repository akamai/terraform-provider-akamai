provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    name = "range_object_rule"
    matches {
      match_type = "clientip"
      object_match_value {
        type  = "range"
        value = [0, 100]
      }
    }
    forward_settings {
      origin_id = "0"
    }
  }
}