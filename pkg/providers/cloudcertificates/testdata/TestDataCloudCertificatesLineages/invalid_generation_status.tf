provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_lineages" "test" {
  generation_status = ["NOT_A_STATUS"]
}
