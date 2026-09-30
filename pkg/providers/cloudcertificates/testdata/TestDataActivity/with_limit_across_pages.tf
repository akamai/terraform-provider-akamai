provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_activity" "test" {
  lineage_id = 500005
  limit      = 150
}
