provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_lineage" "test" {
  contract_id    = "C-0N7RAC7"
  group_id       = 12345
  secure_network = "INVALID_SECURE_NETWORK"
  key_specs = {
    RSA = "2048"
  }
  sans = ["www.example.com"]
}
