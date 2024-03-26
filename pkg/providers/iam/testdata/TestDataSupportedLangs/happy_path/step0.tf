provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_supported_langs" "test" {}
