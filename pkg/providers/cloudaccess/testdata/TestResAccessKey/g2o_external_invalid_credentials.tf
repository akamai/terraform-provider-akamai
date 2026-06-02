provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "credentials" {
  input = {
    key_id = "toolongkey"
    secret = "tooshort"
  }
}

resource "akamai_cloudaccess_key" "test" {
  access_key_name       = "test_key_name"
  authentication_method = "G2O"
  contract_id           = "1-CTRACT"
  credentials_a = {
    cloud_access_key_id     = terraform_data.credentials.output.key_id
    cloud_secret_access_key = terraform_data.credentials.output.secret
    primary_key             = true
  }
  group_id = 12345
  network_configuration = {
    security_network = "ENHANCED_TLS"
  }
}
