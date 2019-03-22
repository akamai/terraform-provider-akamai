provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "default"
  cps_section = "cps"
}

resource "akamai_cps_enrollment" "enrollmenttwo" {
  contract_id = "C-1FRYVV3"
  certificate_chain_type = "default"
  certificate_type = "san"
  change_management = false
  enable_multi_stacked_certificates = false
  ra = "lets-encrypt"
  signature_algorithm = "SHA-256"
  validation_type = "dv"

  admin_contact {
    first_name = "Erik"
    last_name = "Nygren"
    phone = "+33 629586009"
    email = "akamai-devhack@erik.nygren.org"
    address_line_one = "4 Rue Washington"
    address_line_two = ""
    city = "Paris"
    country = "FR"
    organization_name = "Akamai Technologies"
    postal_code = "75008"
    region = "Ile de France"
  }

  csr {
    cn = "testenrollment.akadev.com"
    c = "FR"
    st = "Ile de France"
    l = "Paris"
    o = "Akamai Technologies"
    ou = ""
    sans = ["testenrollment.akadev.com"]
  }

  network_configuration {
    geography = "core"
    must_have_ciphers = "ak-akamai-default-2017q3"
    network_type = "standard-worldwide"
    preferred_ciphers = "ak-akamai-default-2017q3"
    quic_enabled = false
    secure_network = "standard-tls"
    sni_only = true

    dns_name_settings {
      clone_dns_names = true
      dns_names = ["testenrollment.akadev.com"]
    }
  }

  org {
    name = "Akamai Technologies"
    address_line_one = "4 Rue Washington"
    address_line_two = ""
    city = "Paris"
    region = "Ile de France"
    postal_code = "75008"
    country = "FR"
    phone = "+33 1 5669 7161"
  }

  tech_contact {
    first_name = "Erik"
    last_name = "Nygren"
    phone = "+33 1 5669 7161"
    email = "nygren@akamai.com"
    address_line_one = "150 Broadway"
    address_line_two = ""
    city = "Cambridge"
    country = "US"
    organization_name = "Akamai Technologies"
    postal_code = "02142"
    region = "Massachusetts"
  }

  third_party {
    exclude_sans = false
  }
}
