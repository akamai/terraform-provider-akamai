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

data "akamai_apr_user_allow_list" "test" {
  config_id = 96033
}

output "test" {
  value = data.akamai_apr_user_allow_list.test
}
