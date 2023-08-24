provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_imaging_policy_video" "policy" {
  policy {
  }
}