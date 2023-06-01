provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_datacenter" "test" {
  datacenter_id = 1
}
