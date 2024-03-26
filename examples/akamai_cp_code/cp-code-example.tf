terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
}

provider "akamai" {
  edgerc         = "~/.edgerc"
  config_section = "default"
}

resource "akamai_cp_code" "example" {
  contract_id = "ctr_XXX"
  group_id    = "grp_XXX"
  name        = "example-XXX"
  product_id  = "prd_XXX"
}
