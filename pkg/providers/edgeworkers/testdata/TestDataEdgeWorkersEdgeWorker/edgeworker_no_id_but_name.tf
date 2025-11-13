provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworker" "test" {
  name         = "Test Name"
  local_bundle = "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/no_versions.tgz"
}
