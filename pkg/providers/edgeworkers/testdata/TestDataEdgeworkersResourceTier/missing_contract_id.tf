provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_edgeworkers_resource_tier" "test" {
  resource_tier_name = "Basic Compute"
}