provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudaccess_key" "test" {
  access_key_name       = "test_key_name"
  authentication_method = "G2O"
  contract_id           = "1-CTRACT"
  credentials_a = {
    cloud_access_key_id     = "123456789"
    cloud_secret_access_key = "12345678901234567890123456789012"
    primary_key             = true
  }
  group_id = 12345
  network_configuration = {
    security_network = "ENHANCED_TLS"
    additional_cdn   = "CHINA_CDN"
  }
}
