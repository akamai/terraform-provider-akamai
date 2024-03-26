provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudwrapper_activation" "act" {
  config_id = 123
}