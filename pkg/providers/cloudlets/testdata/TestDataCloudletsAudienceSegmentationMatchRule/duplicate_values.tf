provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    matches {
      match_type     = "header"
      match_value    = "this is value"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["no, this is value"]
      }
    }
    forward_settings {
    }
  }
}