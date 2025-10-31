provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_hostname_bindings" "test" {
  network = "foo"
}