provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_generation" "test" {
  lineage_id = 500005
}
