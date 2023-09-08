provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudwrapper_activation" "act" {
  config_id = 123
  revision  = "8b92934d68d69621153c"
  timeouts {
    create = "2s"
    update = "1s"
  }
}