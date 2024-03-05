terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
}

provider "akamai" {}

data "akamai_gtm_default_datacenter" "default_datacenter" {
  domain     = akamai_gtm_domain.tfexample_domain.name
  datacenter = 5400
  depends_on = [
    akamai_gtm_domain.tfexample_domain
  ]
}

// Note: Resource imports must be done one at a time from the command line. Imports can be done in a group by using the Akamai CLI-DNS
// Akamai CLI-DNS will generate import lines, including the Id, for each resource.
// Resource attributes and values should be added to this configuration following import to sync the configuration and state.

resource "akamai_gtm_domain" "tfexample_domain" {
  name = "tfexample.akadns.net"
  type = "weighted"
}

resource "akamai_gtm_datacenter" "tfexample_dc_1" {
  domain = akamai_gtm_domain.tfexample_domain.name
}

resource "akamai_gtm_datacenter" "tfexample_dc_2" {
  domain = akamai_gtm_domain.tfexample_domain.name
}

resource "akamai_gtm_property" "tfexample_prop_1" {
  domain                 = akamai_gtm_domain.tfexample_domain.name
  name                   = "tfexample_prop_1"
  type                   = "weighted-round-robin"
  handout_limit          = 5
  handout_mode           = "normal"
  score_aggregation_type = "median"
}

resource "akamai_gtm_resource" "tfexample_resource_1" {
  domain           = akamai_gtm_domain.tfexample_domain.name
  name             = "tfexample_resource_1"
  aggregation_type = "latest"
  type             = "XML load object via HTTP"
}

resource "akamai_gtm_cidrmap" "tfexample_cidr_1" {
  domain = akamai_gtm_domain.tfexample_domain.name
  name   = "tfexample_cidr_1"
  default_datacenter {
    datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter.datacenter_id
    nickname      = data.akamai_gtm_default_datacenter.default_datacenter.nickname
  }
}

resource "akamai_gtm_asmap" "tfexample_as_1" {
  domain = akamai_gtm_domain.tfexample_domain.name
  name   = "tfexample_as_1"
  default_datacenter {
    datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter.datacenter_id
    nickname      = data.akamai_gtm_default_datacenter.default_datacenter.nickname
  }
}

resource "akamai_gtm_geomap" "tfexample_geo_1" {
  domain = akamai_gtm_domain.tfexample_domain.name
  name   = "tfexample_geo_1"
  default_datacenter {
    datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter.datacenter_id
    nickname      = data.akamai_gtm_default_datacenter.default_datacenter.nickname
  }
}
