provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_api_client" "test" {
  client_id = "c1ien41d"
}