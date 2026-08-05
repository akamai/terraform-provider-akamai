provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_lineage" "test" {
  contract_id    = "C-0N7RAC7"
  group_id       = 12345
  secure_network = "ENHANCED_TLS"
  sans           = ["www.example.com"]
}
