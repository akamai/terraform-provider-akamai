provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudwrapper_activation" "act" {
  config_id = 123
}