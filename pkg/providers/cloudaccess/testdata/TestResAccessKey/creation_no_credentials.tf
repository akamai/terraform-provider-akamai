provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudaccess_key" "test" {
  access_key_name       = "test_key_name"
  authentication_method = "AWS4_HMAC_SHA256"
  contract_id           = "1-CTRACT"
  group_id              = 12345
  network_configuration = {
    security_network = "ENHANCED_TLS"
    additional_cdn   = "CHINA_CDN"
  }
}