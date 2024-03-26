provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_datastreams" "test" {}