provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

variable "create_resource" {
  type    = bool
  default = true
}

variable "access_key_name" {
  type    = string
  default = "test_key_name"
}

variable "authentication_method" {
  type    = string
  default = "AWS4_HMAC_SHA256"
}

variable "contract_id" {
  type    = string
  default = "1-CTRACT"
}

variable "group_id" {
  type    = number
  default = 12345
}

variable "credentials" {
  type = object({
    primary = optional(object({
      cloud_access_key_id     = string
      cloud_secret_access_key = string
      primary_key             = optional(bool, true)
    }))
    secondary = optional(object({
      cloud_access_key_id     = string
      cloud_secret_access_key = string
      primary_key             = optional(bool, false)
    }))
  })
  default = {
    primary = {
      cloud_access_key_id     = "test_key_id"
      cloud_secret_access_key = "test_secret"
      primary_key             = true
    }
    secondary = null
  }
}

variable "network" {
  type = object({
    network_configuration = object({
      security_network = string
      additional_cdn   = optional(string)
    })
  })
  default = {
    network_configuration = {
      security_network = "ENHANCED_TLS"
      additional_cdn   = "CHINA_CDN"
    }
  }
}

resource "akamai_cloudaccess_key" "test" {
  count = var.create_resource ? 1 : 0

  access_key_name       = var.access_key_name
  authentication_method = var.authentication_method
  contract_id           = var.contract_id
  group_id              = var.group_id

  credentials_a = var.credentials.primary
  credentials_b = var.credentials.secondary

  network_configuration = var.network.network_configuration
}