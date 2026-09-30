provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_archived_generations" "test" {
  lineage_id = 500005
}
