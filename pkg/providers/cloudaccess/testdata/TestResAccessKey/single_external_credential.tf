provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

variable "credentials" {
  type = object({
    primary = optional(object({
      cloud_access_key_id     = string
      cloud_secret_access_key = string
      primary_key             = optional(bool, true)
    }))
  })
  default = {
    primary = {
      cloud_access_key_id     = "test_key_id"
      cloud_secret_access_key = "test_secret"
      primary_key             = true
    }
  }
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

  credentials_a = var.credentials.primary
  credentials_b = null
}
