provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_third_party_enrollment" "third_party" {
  contract_id    = "ctr_1"
  common_name    = "test.akamai.com"
  secure_network = "enhanced-tls"
  sni_only       = true
  auto_approve_warnings = [
    "CSR_EXPIRED",
    "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS",
    "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"
  ]
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
  csr {
    country_code        = "US"
    city                = "Cambridge"
    organization        = "Akamai"
    organizational_unit = "WebEx"
    state               = "MA"
  }
  network_configuration {
    geography = "core"
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
