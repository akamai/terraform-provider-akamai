provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgeworkers_resource_tier" "test" {
  resource_tier_name = "Basic Compute"
}