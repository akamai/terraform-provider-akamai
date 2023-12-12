provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_imaging_policy_set" "test_image_set" {
  name        = "test_policy_set"
  region      = "EMEA"
  type        = "IMAGE"
  contract_id = "ctr_1-TEST"
}