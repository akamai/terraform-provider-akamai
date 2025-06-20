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

resource "akamai_apr_user_allow_list" "test" {
  config_id = 96033
  user_allow_list = jsonencode(
    {
      "userAllowListId" : "163481_950K"
    }
  )
}