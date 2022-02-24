provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_imaging_policy_set" "imv_set" {
  name        = "test_policy_set"
  region      = "EMEA"
  type        = "IMAGE"
  contract_id = "1-TEST"
}