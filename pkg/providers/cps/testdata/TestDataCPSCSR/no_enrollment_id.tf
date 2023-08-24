provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cps_csr" "test" {}