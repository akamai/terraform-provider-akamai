terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
    template = {
      source  = "hashicorp/template"
      version = ">= 2.2.0"
    }
  }
}

provider "akamai" {}

data "akamai_group" "group" {
  group_name  = "IPQA Akamai Ion-Express-3-WNKA7W"
  contract_id = "ctr_3-WNKA7W"
}

data "akamai_contract" "contract" {
  group_name = data.akamai_group.group.name
}

data "template_file" "rule_template" {
  template = file("${path.module}/rules/rules.json")
  vars = {
    snippets = "${path.module}/rules/snippets"
  }
}

data "template_file" "rules" {
  for_each = var.customers

  template = data.template_file.rule_template.rendered
  vars = {
    username = each.value.username
    password = each.value.password
  }
}

resource "akamai_cp_code" "cpcode" {

  for_each = var.customers

  product_id  = "prd_Site_Accel"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  name        = each.key
}

resource "akamai_edge_hostname" "edge_hostname" {

  product_id    = "prd_Site_Accel"
  contract_id   = data.akamai_contract.contract.id
  group_id      = data.akamai_group.group.id
  edge_hostname = "test.wheep.co.uk.edgesuite.net"
  ip_behavior   = "IPV6_COMPLIANCE"
}

resource "akamai_property" "property" {

  for_each = var.customers

  name        = each.key
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  product_id  = "prd_Site_Accel"
  rule_format = "v2018-02-27"

  hostnames {
    cname_from             = each.key
    cname_to               = akamai_edge_hostname.edge_hostname.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }
  rules = data.template_file.rules[each.key].rendered
}

resource "akamai_property_activation" "activation" {
  for_each = var.customers

  property_id = akamai_property.property[each.key].id
  contact     = ["you@example.com"]
  network     = upper(var.env)
  version     = akamai_property.property[each.key].latest_version
}
