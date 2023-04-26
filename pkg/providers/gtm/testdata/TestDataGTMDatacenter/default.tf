provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_datacenter" "test" {
  domain        = "test.domain.com"
  datacenter_id = 1
}
