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
