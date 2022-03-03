provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    name = "simple_object_rule"
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      object_match_value {
        type  = "simple"
        value = ["acdc"]
      }
    }
    forward_settings {
      use_incoming_query_string = true
    }
  }
}