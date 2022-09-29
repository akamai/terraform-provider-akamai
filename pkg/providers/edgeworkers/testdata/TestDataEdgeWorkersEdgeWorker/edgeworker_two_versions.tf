provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker" "test" {
  edgeworker_id = 2
  local_bundle  = "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_two_versions.tgz"
}
