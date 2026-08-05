provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_lineages" "test" {
  key_type = "INVALID_KEY_TYPE"
}
