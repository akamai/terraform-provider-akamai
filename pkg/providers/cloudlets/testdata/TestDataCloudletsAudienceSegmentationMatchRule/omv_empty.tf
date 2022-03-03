provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    name = "empty_object_rule"
    matches {
      match_type  = "clientip"
      match_value = "127.0.0.1"
    }
    forward_settings {
    }
  }
}