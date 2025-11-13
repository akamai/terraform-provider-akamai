provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworker" "test" {
  local_bundle = "test_tmp/no_edgeworker_id.tgz"
}
