provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cps_dv_enrollment" "dv" {
  contract_id = "ctr_1"
  common_name = "test.akamai.com"
  secure_network = "enhanced-tls"
  sni_only = true
  acknowledge_pre_verification_warnings = true
  admin_contact {
    first_name = "R1"
    last_name = "D1"
    phone = "123123123"
    email = "r1d1@akamai.com"
    address_line_one = "150 Broadway"
    city = "Cambridge"
    country_code = "US"
    organization = "Akamai"
    postal_code = "12345"
    region = "MA"
  }
  tech_contact {
    first_name = "R2"
    last_name = "D2"
    phone = "123123123"
    email = "r2d2@akamai.com"
    address_line_one = "150 Broadway"
    city = "Cambridge"
    country_code = "US"
    organization = "Akamai"
    postal_code = "12345"
    region = "MA"
  }
  csr {
    country_code = "US"
    city = "Cambridge"
    organization = "Akamai"
    organizational_unit = "WebEx"
    state = "MA"
  }
  network_configuration {
    geography = "core"
  }
  signature_algorithm = "SHA-256"
  organization {
    name = "Akamai"
    phone = "321321321"
    address_line_one = "150 Broadway"
    city = "Cambridge"
    country_code = "US"
    postal_code = "12345"
    region = "MA"
  }
}
