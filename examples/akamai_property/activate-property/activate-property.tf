terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {}

resource "akamai_property_activation" "activation" {
  property = "example.com"
  contact  = ["you@example.com"]
}
