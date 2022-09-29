provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker" "test" {
  edgeworker_id = 1
  local_bundle  = "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_one_version.tgz"
}
