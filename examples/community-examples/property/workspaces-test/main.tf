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

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

variable "edgerc_path" {
  type        = string
  default     = "~/.edgerc"
  description = "Path to edgerc file"
}

variable "config_section" {
  type        = string
  default     = "papi"
  description = "PAPI section"
}

variable "ns_download_domain" {
  type        = string
  description = "Origin hostname"
}

variable "ns_cpcode_id" {
  type        = number
  description = "Origin CP Code"
}

variable "cpcode_name" {
  type        = string
  description = "CP Code name"
}

variable "product" {
  type        = string
  description = "Product"
}

variable "hostname" {
  type        = string
  description = "Hostname"
}

variable "edge_hostname" {
  type        = string
  description = "Edge hostname"
}

variable "email" {
  type        = list(string)
  description = "Notification email"
}

variable "group_name" {
  type        = string
  description = "Access Control Group in which the Akamai configuration should be created"
}

variable "conf_name" {
  type        = string
  description = "Name of the Akamai configuration file"
}

variable "rule_format" {
  type        = string
  description = "PAPI rule schema version"
}

data "akamai_contract" "default" {
  group_name = data.akamai_group.default.name
}

data "akamai_group" "default" {
  contract_id = "test_contract"
  group_name  = var.group_name
}

data "akamai_cp_code" "default" {
  contract_id = data.akamai_contract.default.id
  name        = var.cpcode_name
  group_id    = data.akamai_group.default.id
}

resource "akamai_edge_hostname" "default" {
  product_id    = "prd_${var.product}"
  contract_id   = data.akamai_contract.default.id
  group_id      = data.akamai_group.default.id
  edge_hostname = var.edge_hostname
  ip_behavior   = "IPV6_COMPLIANCE"
}

# Two-stage template evaluation:
# 1. Assemble the snippets using templatefile()
# 2. Interpolate the variables into the resulting template
# Why two-stage: because terraform does not allow recursive
# calls to templatefile.
data "template_file" "rules" {
  # Inner evaluation of the template using templatefile.
  # Purpose: assemble the snippets.
  template = templatefile("${path.module}/rules/rules.tfjson", {
    template_path = "${path.module}/rules"
  })

  vars = {
    ns_download_domain = var.ns_download_domain
    ns_cpcode_id       = var.ns_cpcode_id
    cpcode_name        = var.cpcode_name
    cpcode             = data.akamai_cp_code.default.id
    product            = var.product
  }
}

resource "akamai_property" "default" {
  name        = var.conf_name
  product_id  = "prd_${var.product}"
  contract_id = data.akamai_contract.default.id
  group_id    = data.akamai_group.default.id

  hostnames {
    cname_from             = var.hostname
    cname_to               = akamai_edge_hostname.default.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }

  rule_format = var.rule_format
  rules       = data.template_file.rules.rendered
}

resource "akamai_property_activation" "staging" {
  property_id = akamai_property.default.id
  network     = "STAGING"
  contact     = var.email
  version     = akamai_property.default.latest_version
}

resource "akamai_property_activation" "production" {
  property_id = akamai_property.default.id
  network     = "PRODUCTION"
  contact     = var.email
  version     = akamai_property.default.latest_version
}
