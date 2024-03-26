provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_audience_segmentation_match_rule" "test" {
  match_rules {
    name = "object_object_rule"
    matches {
      match_type = "cookie"
      object_match_value {
        type                = "object"
        name                = "object_match_value_name"
        name_case_sensitive = true
        name_has_wildcard   = false
        options {
          value                = ["object_match_value_value"]
          value_has_wildcard   = false
          value_case_sensitive = true
          value_escaped        = false
        }
      }
    }
    forward_settings {
      path_and_qs = "somePath"
    }
  }
}