terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 0.11.0"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_apr_general_settings" "test" {
  config_id          = 96033
  security_policy_id = "UCON_161669"
}

output "test" {
  value = data.akamai_apr_general_settings.test
}