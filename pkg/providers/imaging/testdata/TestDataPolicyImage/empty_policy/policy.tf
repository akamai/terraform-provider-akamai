provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_imaging_policy_image" "policy" {
  policy {}
}