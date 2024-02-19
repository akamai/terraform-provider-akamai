provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudlets_policy" "policy" {
  name          = "test_policy"
  cloudlet_code = "ER"
  group_id      = "grp_123"
  is_shared     = true
}