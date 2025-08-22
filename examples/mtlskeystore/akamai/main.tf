# This example presents a sample workflow for an `AKAMAI` client certificate. It creates a basic Akamai-signed client certificate, along with
# a CP code, edge hostname, property, and rules that use the `mtls_origin_keystore` behavior to enforce the mTLS Keystore configuration.
# Then, the property is activated on the `PRODUCTION` environment.
#
# To run this example:
#
# 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
#
# The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
#
# 2. Make changes to the attribute values according to your needs.
#
# 3. Open a Terminal or shell instance and initialize the provider with `terraform init`. Then, run `terraform plan` to preview the changes and `terraform apply` to apply your changes.
#
# A successful operation creates an Akamai-signed client certificate, CP code, edge hostname, property, rules, and property activation.

terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 8.1.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = "~/.edgerc"
  config_section = "default"
}

resource "akamai_mtlskeystore_client_certificate_akamai" "akamai_cert" {
  certificate_name    = "Certificate Name"
  contract_id         = "C-0N7RAC7"
  group_id            = 123
  geography           = "CORE"
  notification_emails = ["no-mail@akamai.com"]
  secure_network      = "STANDARD_TLS"
}

resource "akamai_cp_code" "cp_code" {
  contract_id = "C-0N7RAC7"
  group_id    = 123
  product_id  = "prd_SPM"
  name        = "CP-Code-Name"
}

resource "akamai_edge_hostname" "hostname" {
  product_id    = "prd_SPM"
  contract_id   = "C-0N7RAC7"
  group_id      = 123
  ip_behavior   = "IPV6_COMPLIANCE"
  edge_hostname = "www.test-hostname.example.com"
}

resource "akamai_property" "property" {
  name        = "Property-Name"
  contract_id = "C-0N7RAC7"
  group_id    = 123
  product_id  = "prd_SPM"
  hostnames {
    cname_from             = "www.test-hostname-from.example.com"
    cname_to               = akamai_edge_hostname.hostname.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }
  rule_format = data.akamai_property_rules_builder.rule_default.rule_format
  rules       = data.akamai_property_rules_builder.rule_default.json
}

resource "akamai_property_activation" "activation_production" {
  property_id                    = akamai_property.property.id
  contact                        = ["nomail@nomail-akamai.com"]
  version                        = 1
  network                        = "PRODUCTION"
  auto_acknowledge_rule_warnings = true
}

data "akamai_property_rules_builder" "rule_default" {
  rules_v2025_04_29 {
    name      = "default"
    is_secure = false
    behavior {
      origin {
        cache_key_hostname            = "ORIGIN_HOSTNAME"
        compress                      = true
        enable_true_client_ip         = true
        forward_host_header           = "REQUEST_HOST_HEADER"
        hostname                      = "origin-www.example.com"
        http_port                     = 80
        https_port                    = 443
        ip_version                    = "IPV4"
        origin_certificate            = ""
        origin_sni                    = true
        origin_type                   = "CUSTOMER"
        ports                         = ""
        true_client_ip_client_setting = false
        true_client_ip_header         = "True-Client-IP"
        verification_mode             = "PLATFORM_SETTINGS"
      }
    }
    behavior {
      caching {
        behavior        = "MAX_AGE"
        must_revalidate = false
        ttl             = "0s"
      }
    }
    behavior {
      cp_code {
        value {
          id = akamai_cp_code.cp_code.id
        }
      }
    }
    children = [
      data.akamai_property_rules_builder.keystore.json,
    ]
  }
}

data "akamai_property_rules_builder" "keystore" {
  rules_v2025_04_29 {
    name                  = "keystore"
    criteria_must_satisfy = "all"
    behavior {
      mtls_origin_keystore {
        auth_client_cert                = false
        client_certificate_version_guid = akamai_mtlskeystore_client_certificate_akamai.akamai_cert.current_guid
        enable                          = true
      }
    }
  }
}