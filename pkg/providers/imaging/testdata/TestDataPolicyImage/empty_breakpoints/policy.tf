provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_imaging_policy_image" "policy" {
  policy {
    breakpoints {}
  }
}