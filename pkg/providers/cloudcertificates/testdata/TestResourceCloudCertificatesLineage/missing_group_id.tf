provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_lineage" "test" {
  contract_id    = "C-0N7RAC7"
  secure_network = "ENHANCED_TLS"
  key_specs = {
    RSA = "2048"
  }
  sans = ["www.example.com"]
}
