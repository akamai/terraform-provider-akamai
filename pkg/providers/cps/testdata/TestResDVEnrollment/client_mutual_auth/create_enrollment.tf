provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_dv_enrollment" "dv" {
  contract_id    = "ctr_1"
  common_name    = "test.akamai.com"
  secure_network = "enhanced-tls"
  sni_only       = true
  admin_contact {
    first_name       = "R1"
    last_name        = "D1"
    phone            = "123123123"
    email            = "r1d1@akamai.com"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    organization     = "Akamai"
    postal_code      = "12345"
    region           = "MA"
  }
  tech_contact {
    first_name       = "R2"
    last_name        = "D2"
    phone            = "123123123"
    email            = "r2d2@akamai.com"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    organization     = "Akamai"
    postal_code      = "12345"
    region           = "MA"
  }
  certificate_chain_type = "default"
  csr {
    country_code        = "US"
    city                = "Cambridge"
    organization        = "Akamai"
    organizational_unit = "WebEx"
    state               = "MA"
  }
  network_configuration {
    disallowed_tls_versions = [
      "TLSv1",
    "TLSv1_1"]
    clone_dns_names   = false
    geography         = "core"
    ocsp_stapling     = "on"
    preferred_ciphers = "ak-akamai-default"
    must_have_ciphers = "ak-akamai-default"
    quic_enabled      = false
    client_mutual_authentication {
      ocsp_enabled           = true
      send_ca_list_to_client = false
      set_id                 = "12345"
    }
  }
  signature_algorithm = "SHA-256"
  organization {
    name             = "Akamai"
    phone            = "321321321"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    postal_code      = "12345"
    region           = "MA"
  }
}

locals {
  domains_to_validate = [for dvo in akamai_cps_dv_enrollment.dv.dns_challenges : dvo.full_path]
}

output "domains_to_validate" {
  value = join(",", local.domains_to_validate)
}
