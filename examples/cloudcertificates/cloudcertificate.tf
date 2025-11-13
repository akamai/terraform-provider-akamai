# This example presents a sample CCM workflow that includes creating a self-signed cloud certificate and uploading it.
# Optionally, provide a PEM-encoded trust chain when uploading the signed certificate.
#
# Before applying this example, make changes to the attribute values according to your needs.
#
# This workflow generates a self-signed certificate, provisions a cloud certificate, and uploads the signed certificate.
#
# Use the certificate ID from the `akamai_cloudcertificates_upload_signed_certificate` resource to bind the signed certificate with the hostname in the `akamai_property` resource.

terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.2.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = "~/.edgerc"
  config_section = "default"
}

resource "akamai_cloudcertificates_certificate" "test" {
  base_name      = "example-base-name"
  contract_id    = "C-0N7RAC7"
  group_id       = "grp_123"
  key_size       = "2048"
  key_type       = "RSA"
  secure_network = "ENHANCED_TLS"
  sans           = ["test.example.com", "www.test2.example.com"]
  subject = {
    common_name  = "test.example.com"
    organization = "Test Organization"
    state        = "Massachusetts"
    locality     = "Cambridge"
    country      = "US"
  }
}

resource "tls_private_key" "key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "cert" {
  depends_on            = [tls_private_key.key]
  private_key_pem       = tls_private_key.key.private_key_pem
  validity_period_hours = 8760
  is_ca_certificate     = true
  allowed_uses = [
    "cert_signing",
    "key_encipherment",
    "digital_signature",
    "crl_signing"
  ]
  subject {
    common_name  = "test.com"
    organization = "Akamai"
  }
}

resource "tls_locally_signed_cert" "signed_cert" {
  depends_on            = [tls_private_key.key]
  ca_private_key_pem    = tls_private_key.key.private_key_pem
  cert_request_pem      = akamai_cloudcertificates_certificate.test.csr_pem
  ca_cert_pem           = tls_self_signed_cert.cert.cert_pem
  validity_period_hours = 8760
  allowed_uses = [
    "cert_signing",
    "key_encipherment",
    "digital_signature",
    "crl_signing"
  ]
}

resource "akamai_cloudcertificates_upload_signed_certificate" "upload" {
  certificate_id         = akamai_cloudcertificates_certificate.test.certificate_id
  acknowledge_warnings   = true
  signed_certificate_pem = tls_locally_signed_cert.signed_cert.cert_pem
}