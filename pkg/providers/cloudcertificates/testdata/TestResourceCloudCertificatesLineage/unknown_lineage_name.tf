provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "lineage_name" {
  input = "unknown-lineage-name"
}

resource "akamai_cloudcertificates_lineage" "test" {
  contract_id    = "C-0N7RAC7"
  group_id       = 12345
  secure_network = "ENHANCED_TLS"
  lineage_name   = terraform_data.lineage_name.output
  key_specs = {
    RSA = "2048"
  }
  sans = ["www.example.com", "example.com"]
}
