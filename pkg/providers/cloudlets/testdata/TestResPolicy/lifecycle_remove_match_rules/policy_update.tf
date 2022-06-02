provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudlets_policy" "policy" {
  name          = "test_policy"
  cloudlet_code = "ER"
  description   = "test policy description"
  group_id      = "grp_123"
}