terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source = "registry.terraform.io/akamai/akamai"
    }
  }
}

data "akamai_group" "group" {
  group_name  = "Group-1"
  contract_id = "Contract-1"
}

resource "akamai_apidefinitions_api" "api" {
  api         = file("${path.module}/api.json")
  contract_id = trimprefix(data.akamai_group.group.contract_id, "ctr_")
  group_id    = trimprefix(data.akamai_group.group.id, "grp_")
}

resource "akamai_apidefinitions_activation" "api_activation_staging" {
  api_id                    = akamai_apidefinitions_api.api.id
  version                   = akamai_apidefinitions_api.api.latest_version
  network                   = "STAGING"
  notification_recipients   = ["user@example.com"]
  notes                     = "Notes"
  auto_acknowledge_warnings = true
}

resource "akamai_apidefinitions_activation" "api_activation_production" {
  api_id                    = akamai_apidefinitions_api.api.id
  version                   = akamai_apidefinitions_api.api.latest_version
  network                   = "PRODUCTION"
  notification_recipients   = ["user@example.com"]
  notes                     = "Notes"
  auto_acknowledge_warnings = true
}