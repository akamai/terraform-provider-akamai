provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_imaging_policy_set" "test_image_set" {
  name        = "test_policy_set"
  region      = "EMEA"
  type        = "IMAGE"
  contract_id = "1-TEST"
}