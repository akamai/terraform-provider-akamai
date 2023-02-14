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

resource "akamai_property" "property" {
  name = "akavadeveloper"

  product  = "prd_SPM"
  cp_code  = akamai_cp_code.cp_code.id
  contract = data.akamai_contract.contract.id
  group    = data.akamai_group.group.id

  hostnames = {
    "terraform.example.org" = akamai_edge_hostname.test.edge_hostname
  }

  rule_format = "v2019-07-25"

  rules = data.akamai_property_rules.rules.json
}

data "akamai_contract" "contract" {}

data "akamai_group" "group" {}

resource "akamai_cp_code" "cp_code" {
  name     = "terraform-testing"
  contract = data.akamai_contract.contract.id
  group    = data.akamai_group.group.id
  product  = "prd_SPM"
}

resource "akamai_edge_hostname" "test" {
  product       = "prd_SPM"
  contract      = data.akamai_contract.contract.id
  group         = data.akamai_group.group.id
  edge_hostname = "terraform-test1.akavadeveloper.io.edgesuite.net"
  ipv4          = true
  ipv6          = true
}

data "akamai_property_rules" "rules" {
  rules {
    behavior {
      name = "origin"
      option {
        key   = "cacheKeyHostname"
        value = "ORIGIN_HOSTNAME"
      }
      option {
        key   = "compress"
        value = "true"
      }
      option {
        key   = "enableTrueClientIp"
        value = "false"
      }
      option {
        key   = "forwardHostHeader"
        value = "REQUEST_HOST_HEADER"
      }
      option {
        key   = "hostname"
        value = "exampleakavadeveloper.io"
      }
      option {
        key   = "httpPort"
        value = "80"
      }
      option {
        key   = "httpsPort"
        value = "443"
      }
      option {
        key   = "originSni"
        value = "true"
      }
      option {
        key   = "originType"
        value = "CUSTOMER"
      }
      option {
        key   = "verificationMode"
        value = "PLATFORM_SETTINGS"
      }
      option {
        key   = "originCertificate"
        value = ""
      }
      option {
        key   = "ports"
        value = ""
      }
    }
    behavior {
      name = "cpCode"
      option {
        key   = "id"
        value = akamai_cp_code.cp_code.id
      }
      option {
        key   = "name"
        value = akamai_cp_code.cp_code.name
      }
    }
    behavior {
      name = "caching"
      option {
        key   = "behavior"
        value = "MAX_AGE"
      }
      option {
        key   = "mustRevalidate"
        value = "false"
      }
      option {
        key   = "ttl"
        value = "1d"
      }
    }
  }
}
