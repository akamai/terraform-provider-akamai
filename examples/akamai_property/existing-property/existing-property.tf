terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
}

provider "akamai" {}

resource "akamai_property" "prop" {
  name        = "host.example.com"
  product_id  = "prd_SPM"
  group_id    = "grp_1"
  contract_id = "ctr_1"
  hostnames {
    cname_from             = "examplehost"
    cname_to               = "host.example.com"
    cert_provisioning_type = "CPS_MANAGED"
  }

  rules = data.akamai_property_rules.prop.rules
}

data "akamai_property_rules" "prop" {
  property_id = "prp_1"
}
