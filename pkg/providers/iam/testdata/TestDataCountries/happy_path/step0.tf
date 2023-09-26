provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_iam_countries" "test" {}
