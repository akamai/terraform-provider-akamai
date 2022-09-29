provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgeworker" "test" {
  local_bundle_path = "test_tmp/no_edgeworker_id.tgz"
}
