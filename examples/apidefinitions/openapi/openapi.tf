terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source = "registry.terraform.io/akamai/akamai"
    }
  }
}

data "akamai_apidefinitions_openapi" "petstore" {
  file = file("${path.module}/petstore-3.0.yml")
}

resource "akamai_apidefinitions_api" "api" {
  api = data.akamai_apidefinitions_openapi.petstore.api
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
