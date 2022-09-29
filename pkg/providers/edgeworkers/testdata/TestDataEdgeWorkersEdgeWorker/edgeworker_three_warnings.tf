provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker" "test" {
  edgeworker_id = 4
  local_bundle  = "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_three_warnings.tgz"
}
