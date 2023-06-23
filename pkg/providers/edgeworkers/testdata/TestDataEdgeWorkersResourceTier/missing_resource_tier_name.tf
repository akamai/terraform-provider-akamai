provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworkers_resource_tier" "test" {
  contract_id = "1-599K"
}