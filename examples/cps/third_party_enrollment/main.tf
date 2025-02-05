# This example presents a sample workflow for a third-party enrollment. It creates a basic third-party enrollment,
# fetches the Certificate Signing Request (CSR) for this enrollment, and then uploads the signed certificate with the `akamai_cps_upload_certificate` resource.
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
# A successful operation creates a third-party enrollment and fetches its CSR (Certificate Signing Request).
#
# 4. Use the RSA or ECDSA CSR from the `akamai_cps_csr` data source to retrieve the signed certificate with a third-party tool.
# Uncomment the `akamai_cps_upload_certificate` resource and use the signed certificate in the `certificate_ecdsa_pem` or `certificate_rsa_pem` attribute.
# Optionally, you can provide the trust chains.
#
# 5. Run `terraform apply`.
#
# A successful operation deploys the certificate to the `STAGING` network. To push the changes to
# the `PRODUCTION` network, change the `acknowledge_change_management` attribute to `true`.

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

resource "akamai_cps_third_party_enrollment" "tp_enrollment" {
  contract_id         = "C-0N7RAC7"
  common_name         = "*.example.com"
  secure_network      = "enhanced-tls"
  sni_only            = true
  change_management   = true
  signature_algorithm = "SHA-256"
  # You can fetch the list of available warning values using the `akamai_cps_warnings` data source.
  auto_approve_warnings = [
    "DNS_NAME_LONGER_THEN_255_CHARS",
    "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS",
    "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"
  ]
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

data "akamai_cps_csr" "my-csr" {
  enrollment_id = akamai_cps_third_party_enrollment.tp_enrollment.id
}

#resource "akamai_cps_upload_certificate" "upload_cert" {
#  enrollment_id                          = akamai_cps_third_party_enrollment.tp_enrollment.id
#  certificate_ecdsa_pem                  = example_cert_ecdsa.pem
#  trust_chain_ecdsa_pem                  = example_trust_chain_ecdsa.pem
#  acknowledge_post_verification_warnings = true
#  ### After testing on the `STAGING` network, change the `acknowledge_change_management` attribute to `true` to deploy the changes
#  ### to the `PRODUCTION` network.
#  acknowledge_change_management          = false
#  wait_for_deployment                    = true
#  timeouts {
#    default = "1h"
#  }
#}