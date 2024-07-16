provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudaccess_key" "test" {
  access_key_name       = "test_key_name"
  authentication_method = "AWS4_HMAC_SHA256"
  contract_id           = "1-CTRACT"
  credentials_a = {
    cloud_access_key_id     = "test_key_id"
    cloud_secret_access_key = "test_secret"
    primary_key             = true
  }
  credentials_b = {
    cloud_access_key_id     = "test_key_id_2"
    cloud_secret_access_key = "test_secret_2"
    primary_key             = false
  }
  group_id = 12345
  network_configuration = {
    security_network = "ENHANCED_TLS"
    additional_cdn   = "CHINA_CDN"
  }
}