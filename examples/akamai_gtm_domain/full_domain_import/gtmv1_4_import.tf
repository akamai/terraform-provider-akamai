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

// Note: Resource imports must be done one at a time from the command line. Imports can be done in a group by using the Akamai CLI-DNS
// Akamai CLI-DNS will generate import lines, including the Id, for each resource.
// Resource attributes and values should be added to this configuration following import to sync the configuration and state.

resource "akamai_gtm_domain" "tfexample_domain" {
}

resource "akamai_gtm_datacenter" "tfexample_dc_1" {
}

resource "akamai_gtm_datacenter" "tfexample_dc_2" {
}

resource "akamai_gtm_property" "tfexample_prop_1" {
}

resource "akamai_gtm_resource" "tfexample_resource_1" {
}

resource "akamai_gtm_cidrmap" "tfexample_cidr_1" {
}

resource "akamai_gtm_asmap" "tfexample_as_1" {
}

resource "akamai_gtm_geomap" "tfexample_geo_1" {
}
