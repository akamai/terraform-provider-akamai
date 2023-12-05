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

resource "akamai_property" "property" {
  name = "akavadeveloper"

  product_id  = "prd_SPM"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id

  hostnames {
    cname_from             = "terraform.example.org"
    cname_to               = akamai_edge_hostname.test.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }

  rule_format = "v2019-07-25"

  rules = data.akamai_property_rules.rules.rules
}

data "akamai_group" "group" {
  group_name  = "Example.com-1-1TJZH5"
  contract_id = "ctr_1-1TJZH5"
}

data "akamai_contract" "contract" {
  group_name = data.akamai_group.group.group_name
}

resource "akamai_cp_code" "cp_code" {
  name        = "terraform-testing"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  product_id  = "prd_SPM"
}

resource "akamai_edge_hostname" "test" {
  product_id    = "prd_SPM"
  contract_id   = data.akamai_contract.contract.id
  group_id      = data.akamai_group.group.id
  edge_hostname = "terraform-test1.akavadeveloper.io.edgesuite.net"
  ip_behavior   = "IPV6_COMPLIANCE"
}

data "akamai_property_rules" "rules" {
  property_id = "prp_1"
}
