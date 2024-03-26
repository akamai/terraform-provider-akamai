provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_imaging_policy_video" "policy" {
  policy {
  }
}