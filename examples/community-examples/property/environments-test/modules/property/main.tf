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

data "akamai_group" "group" {
  group_name  = "IPQA Akamai Ion-Express-3-WNKA7W"
  contract_id = "ctr_3-WNKA7W"
}

data "akamai_contract" "contract" {
  group_name = data.akamai_group.group.name
}

data "template_file" "rules" {
  template = file("${path.module}/rules.json")
}

resource "akamai_cp_code" "test-wheep-co-uk" {
  product_id  = "prd_Download_Delivery"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  name        = "${var.env}.wheep.co.uk"
}

resource "akamai_edge_hostname" "test-wheep-co-uk-edgesuite-net" {
  product_id    = "prd_Download_Delivery"
  contract_id   = data.akamai_contract.contract.id
  group_id      = data.akamai_group.group.id
  ip_behavior   = "IPV6_COMPLIANCE"
  edge_hostname = "${var.env}.wheep.co.uk.edgesuite.net"
}

resource "akamai_property" "test-wheep-co-uk" {
  name        = "${var.env}.wheep.co.uk"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  product_id  = "prd_Download_Delivery"
  rule_format = "latest"
  hostnames {
    cname_from             = "${var.env}.wheep.co.uk"
    cname_to               = akamai_edge_hostname.test-wheep-co-uk-edgesuite-net.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }
  rules = data.template_file.rules.rendered
}

resource "akamai_property_activation" "test-wheep-co-uk" {
  property_id = akamai_property.test-wheep-co-uk.id
  contact     = ["you@example.com"]
  network     = upper(var.network)
  version     = akamai_property.test-wheep-co-uk.latest_version
}
