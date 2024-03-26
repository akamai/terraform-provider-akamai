provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudwrapper_activation" "act" {
  config_id = 321
  revision  = "5fe7963eb7270e69c5e8"
}