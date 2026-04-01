provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudaccess_key" "test" {
  access_key_name       = "test_key_name"
  authentication_method = "AVM_CLOUDINARY"
  contract_id           = "1-CTRACT"
  credentials_b = {
    cloud_secret_access_key = "test_secret"
    primary_key             = true
  }
  credentials_a = {
    cloud_secret_access_key = "test_secret_2"
    primary_key             = false
  }
  group_id = 12345
  network_configuration = {
    security_network = "ENHANCED_TLS"
  }
}
