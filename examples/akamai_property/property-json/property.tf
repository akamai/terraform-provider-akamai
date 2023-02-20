terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
    local = {
      source  = "hashicorp/local"
      version = ">= 1.1.0"
    }
  }
}

provider "akamai" {}

resource "akamai_property" "terraform_example" {
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  name        = "terraform_example1"
  product_id  = "prd_SPM"

  hostnames {
    cname_from             = "terraform.example1.org"
    cname_to               = akamai_edge_hostname.ehn.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }

  hostnames {
    cname_from             = "terraform.example1.com"
    cname_to               = akamai_edge_hostname.ehn.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }

  rule_format = "v2019-07-25"
  rules       = data.local_file.rules.content
}

resource "akamai_edge_hostname" "ehn" {
  edge_hostname = "terraform.example1.org.edgesuite.net"

  product_id  = "prd_SPM"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id

  ip_behavior = "IPV6_COMPLIANCE"
}

data "akamai_contract" "contract" {
  group_name = data.akamai_group.group.name
}

data "akamai_group" "group" {
  group_name  = "Example.com-1-1TJZH5"
  contract_id = "ctr_1-1TJZH5"
}

data "local_file" "rules" {
  filename = "rules.json"
}
