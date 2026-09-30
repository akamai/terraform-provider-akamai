provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_lineage" "test" {
  group_id       = 12345
  secure_network = "ENHANCED_TLS"
  key_specs = {
    RSA = "2048"
  }
  sans = ["www.example.com"]
}
