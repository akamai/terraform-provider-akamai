terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.0.0"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_apr_general_settings" "test" {
  config_id          = 96033
  security_policy_id = "UCON_161669"
  general_settings = jsonencode(
    {
      accountProtection           = true
      originSignalHeader          = false
      originUserIdInRequestHeader = false
      usernameInRequestHeader     = true
    }
  )
}
