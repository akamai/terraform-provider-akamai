provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_imaging_policy_video" "policy" {
  policy {
  }
}