terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
    template = {
      source = "hashicorp/template"
    }
  }
}

provider "akamai" {}

variable "edgerc" {
  type        = string
  default     = "~/.edgerc"
  description = "Path to edgerc file"
}

variable "edgerc_papi" {
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

variable "production" {
  type        = bool
  default     = false
  description = "Deploy to Akamai production network"
}

variable "staging" {
  type        = bool
  default     = true
  description = "Deploy to Akamai staging network"
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

data "akamai_contract" "default" {}

data "akamai_group" "default" {
  contract = data.akamai_contract.default.id
  name     = var.group_name
}

data "akamai_cp_code" "default" {
  contract = data.akamai_contract.default.id
  name     = var.cpcode_name
  group    = data.akamai_group.default.id
}

resource "akamai_edge_hostname" "default" {
  product       = "prd_${var.product}"
  contract      = data.akamai_contract.default.id
  group         = data.akamai_group.default.id
  edge_hostname = var.edge_hostname
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
    cpcode             = parseint(replace(data.akamai_cp_code.default.id, "cpc_", ""), 10)
    product            = var.product
  }
}

resource "akamai_property" "default" {
  name    = var.conf_name
  contact = var.email

  product  = "prd_${var.product}"
  contract = data.akamai_contract.default.id
  group    = data.akamai_group.default.id

  hostnames = {
    (var.hostname) = akamai_edge_hostname.default.edge_hostname
  }

  rule_format = var.rule_format
  rules       = data.template_file.rules.rendered
}

resource "akamai_property_activation" "staging" {
  property = akamai_property.default.id
  network  = "STAGING"
  activate = var.staging
  contact  = var.email
}

resource "akamai_property_activation" "production" {
  property = akamai_property.default.id
  network  = "PRODUCTION"
  activate = var.production
  contact  = var.email
}
