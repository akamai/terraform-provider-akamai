provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_datacenter" "test" {
  domain        = "test.domain.com"
  datacenter_id = 1
}
