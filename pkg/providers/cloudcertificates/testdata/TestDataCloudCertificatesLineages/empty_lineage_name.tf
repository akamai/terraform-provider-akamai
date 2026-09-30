provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_lineages" "test" {
  lineage_name = ""
}
