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

data "akamai_apr_protected_operations" "test" {
  config_id          = 96033
  security_policy_id = "UCON_161669"
  operation_id       = "e93a594a-5798-437b-a5e3-116107a953f7"
}

output "test" {
  value = data.akamai_apr_protected_operations.test
}