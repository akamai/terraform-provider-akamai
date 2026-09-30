provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_lineage" "test" {
  contract_id    = "C-0N7RAC7"
  group_id       = 12345
  lineage_name   = "explicit-lineage"
  lineage_type   = "SINGLE_GENERATION"
  secure_network = "ENHANCED_TLS"
  key_specs = {
    RSA = "2048"
  }
  sans = ["www.example.com", "example.com"]
}
