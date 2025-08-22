# This example presents a sample workflow for a third-party enrollment
# with client mutual authentication enabled by referring to the activated CA set.
# It fetches the certificate signing request (CSR) for this enrollment, uses a self-signed certificate to sign it, 
# and then uploads the signed certificate with the `akamai_cps_upload_certificate` resource.
#
# Before applying this example, make changes to the attribute values according to your needs.
# 
# Make sure to use the activated CA set ID from the `akamai_mtlstruststore_ca_set_activation` resource via the `set_id` attribute, 
# inside the `network_configuration.client_mutual_authentication` block.
#
# A successful operation:
# - Creates a third-party enrollment and fetches its certificate signing request (CSR).
# - Uses the `tls_self_signed_cert` and `tls_locally_signed_cert` resources to generate a self-signed certificate and sign the CSR.
# - Uploads the signed certificate to the CPS with the `akamai_cps_upload_certificate` resource.
# - Deploys the certificate to the `STAGING` and `PRODUCTION` environments.
#
# To seamlessly remove this configuration, you must remove this resource before the `akamai_mtlstruststore_ca_set_activation` resource.
# If you are referring to the `PRODUCTION` activation resource, you can add the `depends_on` block to ensure the correct order of resource destruction,
# making sure that the `STAGING` activation is also removed after the enrollment resource removal.

resource "akamai_cps_third_party_enrollment" "enrollment" {
  depends_on                            = [akamai_mtlstruststore_ca_set_activation.ca_set_activation_staging]
  contract_id                           = "C-0N7RAC7"
  common_name                           = "*.example.com"
  secure_network                        = "enhanced-tls"
  sni_only                              = true
  acknowledge_pre_verification_warnings = true
  change_management                     = false
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
    client_mutual_authentication {
      send_ca_list_to_client = true
      ocsp_enabled           = true
      set_id                 = akamai_mtlstruststore_ca_set_activation.ca_set_activation_production.ca_set_id
    }
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
}

data "akamai_cps_csr" "cps_csr" {
  enrollment_id = akamai_cps_third_party_enrollment.enrollment.id
}

resource "tls_private_key" "cps_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "cps_certificate" {
  private_key_pem       = tls_private_key.cps_key.private_key_pem
  validity_period_hours = 8760
  is_ca_certificate     = true

  allowed_uses = [
    "cert_signing",
    "key_encipherment",
    "digital_signature",
    "crl_signing"
  ]

  subject {
    common_name  = "example.com"
    organization = "Akamai"
  }
}

resource "tls_locally_signed_cert" "signed_certificate" {
  ca_private_key_pem    = tls_private_key.cps_key.private_key_pem
  cert_request_pem      = trimspace(data.akamai_cps_csr.cps_csr.csr_ecdsa)
  ca_cert_pem           = tls_self_signed_cert.cps_certificate.cert_pem
  validity_period_hours = 8760

  allowed_uses = [
    "cert_signing",
    "key_encipherment",
    "digital_signature",
    "crl_signing"
  ]
}

resource "akamai_cps_upload_certificate" "upload_cert" {
  enrollment_id                          = akamai_cps_third_party_enrollment.enrollment.id
  certificate_ecdsa_pem                  = tls_locally_signed_cert.signed_certificate.cert_pem
  acknowledge_post_verification_warnings = true
  wait_for_deployment                    = true
}