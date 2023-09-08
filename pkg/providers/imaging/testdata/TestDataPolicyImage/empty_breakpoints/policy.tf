provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_imaging_policy_image" "policy" {
  policy {
    breakpoints {}
  }
}