provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_resource" "my_gtm_resource" {
  resource_name = "resource1"
}