provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_iam_supported_langs" "test" {}
