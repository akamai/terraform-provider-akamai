terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {}

resource "akamai_cp_code" "example" {
  contract = "ctr_XXX"
  group    = "grp_XXX"
  name     = "example-XXX"
  product  = "prd_XXX"
}
