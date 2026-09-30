provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_lineage" "test" {
  contract_id    = "C-0N7RAC7"
  group_id       = 12345
  geo_class      = "STANDARD_WORLDWIDE"
  secure_network = "ENHANCED_TLS"
  key_specs = {
    RSA   = "2048"
    ECDSA = "P-256"
  }
  sans = ["www.example.com", "example.com"]
  subject = {
    common_name         = "example.com"
    organization        = "Example Corp."
    organizational_unit = "IT"
    country             = "US"
    state               = "Massachusetts"
    locality            = "Cambridge"
  }
}

# stands in for a real CA-signing chain (e.g. tls_private_key/tls_self_signed_cert/tls_locally_signed_cert):
# its output is the signed certificate content that would otherwise come from signing the lineage's CSR.
resource "terraform_data" "rsa_cert" {
  input = "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"
}

resource "akamai_cloudcertificates_upload" "upload" {
  lineage_id = akamai_cloudcertificates_lineage.test.lineage_id
  algorithms = {
    RSA = {
      signed_certificate_pem = terraform_data.rsa_cert.output
    }
  }
}

resource "akamai_cloudcertificates_activation" "activation" {
  lineage_id               = akamai_cloudcertificates_upload.upload.lineage_id
  staging_generation_id    = akamai_cloudcertificates_upload.upload.generation_id
  production_generation_id = akamai_cloudcertificates_upload.upload.generation_id
}
