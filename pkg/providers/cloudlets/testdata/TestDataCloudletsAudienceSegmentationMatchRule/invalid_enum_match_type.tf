provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    matches {
      match_type     = "invalid_match_type"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["a"]
      }
    }
    forward_settings {
    }
  }
}