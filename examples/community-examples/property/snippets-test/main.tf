terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
    template = {
      source  = "hashicorp/template"
      version = "~> 0.1"
    }
  }
}

provider "akamai" {}

data "akamai_group" "group" {}

data "akamai_contract" "contract" {
  group = data.akamai_group.group.name
}

data "template_file" "rule_template" {
  template = file("${path.module}/rules/rules.json")
  vars = {
    snippets = "${path.module}/rules/snippets"
  }
}

data "template_file" "rules" {
  template = data.template_file.rule_template.rendered
  vars = {
    tdenabled = true
  }
}

resource "akamai_cp_code" "test-wheep-co-uk" {
  product  = "prd_Site_Accel"
  contract = data.akamai_contract.contract.id
  group    = data.akamai_group.group.id
  name     = "test-wheep-co-uk"
}

resource "akamai_edge_hostname" "test-wheep-co-uk" {
  product       = "prd_Site_Accel"
  contract      = data.akamai_contract.contract.id
  group         = data.akamai_group.group.id
  edge_hostname = "tf2.wheep.co.uk.edgesuite.net"
}

resource "akamai_property" "test-wheep-co-uk" {
  name        = "tfsnippets.wheep.co.uk"
  cp_code     = akamai_cp_code.test-wheep-co-uk.id
  contract    = data.akamai_contract.contract.id
  group       = data.akamai_group.group.id
  product     = "prd_Site_Accel"
  rule_format = "v2018-02-27"

  hostnames = {
    "tfsnippets.wheep.co.uk"   = akamai_edge_hostname.test-wheep-co-uk.edge_hostname,
    "testsnippets.wheep.co.uk" = akamai_edge_hostname.test-wheep-co-uk.edge_hostname
  }
  rules     = data.template_file.rules.rendered
  is_secure = true

}

resource "akamai_property_activation" "test-wheep-co-uk" {
  property = akamai_property.test-wheep-co-uk.id
  contact  = ["you@example.com"]
  network  = upper(var.env)
  activate = true
}
