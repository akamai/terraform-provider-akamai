provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudlets_policy" "policy" {
  name          = "test_policy"
  cloudlet_code = "ER"
  group_id      = "grp_123"
}