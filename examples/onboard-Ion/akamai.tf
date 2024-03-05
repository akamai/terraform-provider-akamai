terraform {
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
  required_version = ">= 1.0"
}

provider "akamai" {
  #Credentials can be provided inline using services such as Terraform Vault or by via an .edgerc file
  edgerc         = "../../.edgerc"
  config_section = "papi"
  #  property {
  #        host = "${var.akamai_host}"
  #        access_token = "${var.akamai_access_token}"
  #        client_token = "${var.akamai_client_token}"
  #        client_secret = "${var.akamai_client_secret}"
  #    }
}

data "akamai_group" "group" {
  group_name  = "test"
  contract_id = "test_contract"
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
  template = data.template_file.rule_template.rendered
  vars = {
    tdenabled = var.tdenabled
  }
}

resource "akamai_edge_hostname" "example-property" {
  product_id    = "prd_xxxx"
  contract_id   = data.akamai_contract.contract.id
  group_id      = data.akamai_group.group.id
  edge_hostname = "xxxx.edgesuite.net"
  ip_behavior   = "IPV6_COMPLIANCE"
}

resource "akamai_property" "example-property" {
  name        = "example.mydomain.com"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  product_id  = "prd_xxxx"
  rule_format = "latest"
  hostnames {
    cname_from             = "example.mydomain.com"
    cname_to               = akamai_edge_hostname.example-property.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }
  rules = data.template_file.rules.rendered
}

resource "akamai_property_activation" "example-property" {
  property_id = akamai_property.example-property.id
  contact     = ["me@mydomain.com"]
  network     = "STAGING"
  version     = akamai_property.example-property.latest_version
}
