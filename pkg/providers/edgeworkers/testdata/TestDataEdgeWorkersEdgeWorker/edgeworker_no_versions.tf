provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker" "test" {
  edgeworker_id = 5
  local_bundle  = "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/no_versions.tgz"
}
