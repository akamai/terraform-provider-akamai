provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_third_party_enrollment" "test" {
  contract_id                           = "ctr_1"
  acknowledge_pre_verification_warnings = false
  auto_approve_warnings                 = []
  change_management                     = false
  common_name                           = "test.akamai.com"
  secure_network                        = "enhanced-tls"
  sni_only                              = true
  admin_contact {
    first_name       = "R1"
    last_name        = "D1"
    organization     = "Akamai"
    email            = "r1d1@akamai.com"
    phone            = "123123123"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    region           = "MA"
    postal_code      = "12345"
    country_code     = "US"
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
    enable_for_all_sans = true
    dns_names = [
      "test.akamai.com",
    ]
    geography         = "core"
    ocsp_stapling     = "on"
    preferred_ciphers = "ak-akamai-default"
    must_have_ciphers = "ak-akamai-default"
    quic_enabled      = false
  }
  organization {
    name             = "Akamai"
    phone            = "321321321"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    postal_code      = "12345"
    region           = "MA"
  }
  signature_algorithm = "SHA-256"
  tech_contact {
    first_name       = "R2"
    last_name        = "D2"
    organization     = "Akamai"
    email            = "r2d2@akamai.com"
    phone            = "123123123"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    region           = "MA"
    postal_code      = "12345"
    country_code     = "US"
  }
}
