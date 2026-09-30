provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_bindings" "test" {
  lineage_id = 500005
  sort_order = "NOT_A_SORT_ORDER"
}
