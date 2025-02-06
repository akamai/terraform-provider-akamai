provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_api_clients" "test" {
}
