provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cps_warnings" "test" {}
