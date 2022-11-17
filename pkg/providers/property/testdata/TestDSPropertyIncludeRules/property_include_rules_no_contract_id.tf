provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_include_rules" "test" {
  group_id    = "grp_2"
  include_id  = "12345"
  version     = 1
  name        = "TestIncludeName"
  type        = "MICROSERVICES"
  rule_format = "v2022-06-28"
}
