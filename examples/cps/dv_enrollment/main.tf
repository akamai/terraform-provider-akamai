# This example presents a sample workflow for a DV enrollment. It creates a basic DV enrollment and validates
# that certificate, so it can be deployed to `STAGING` and `PRODUCTION` networks.
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
# A successful operation creates a DV enrollment.
#
# 4. Complete the DNS or HTTP challenges for the domain.
#
# 5. Uncomment the `akamai_cps_dv_validation` resource. Then, run `terraform plan` to preview the changes and `terraform apply` to apply your changes.
#
# A successful operation deploys DV certificate to the `STAGING` and `PRODUCTION` networks.

terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 3.1.0"
    }
  }
}

provider "akamai" {
  edgerc         = "~/.edgerc"
  config_section = "default"
}

resource "akamai_cps_dv_enrollment" "dv_enrollment" {
  contract_id                           = "C-0N7RAC7"
  acknowledge_pre_verification_warnings = true
  common_name                           = "*.example.com"
  secure_network                        = "enhanced-tls"
  sni_only                              = true
  signature_algorithm                   = "SHA-256"
  admin_contact {
    first_name       = "John"
    last_name        = "Smith"
    phone            = "+1-311-555-2368"
    email            = "jsmith@example.com"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    organization     = "Example Corp."
    postal_code      = "02142"
    region           = "MA"
    title            = "Administrator"
  }
  tech_contact {
    first_name       = "Jane"
    last_name        = "Smith"
    phone            = "+1-311-555-2369"
    email            = "jasmith@akamai.com"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    organization     = "Example Corp."
    postal_code      = "02142"
    region           = "MA"
    title            = "Administrator"
  }
  certificate_chain_type = "default"
  csr {
    country_code        = "US"
    city                = "Cambridge"
    organization        = "Example Corp."
    organizational_unit = "Corp IT"
    state               = "MA"
  }
  network_configuration {
    disallowed_tls_versions = ["TLSv1", "TLSv1_1"]
    clone_dns_names         = false
    geography               = "core"
    ocsp_stapling           = "on"
    preferred_ciphers       = "ak-akamai-2020q1"
    must_have_ciphers       = "ak-akamai-2020q1"
    quic_enabled            = false
  }
  organization {
    name             = "Example Corp."
    phone            = "+1-311-555-2370"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    postal_code      = "02142"
    region           = "MA"
  }
  timeouts {
    default = "1h"
  }
}

#resource "akamai_cps_dv_validation" "example" {
#  enrollment_id                          = akamai_cps_dv_enrollment.dv_enrollment.id
#  acknowledge_post_verification_warnings = true
#  timeouts {
#    default = "1h"
#  }
#}