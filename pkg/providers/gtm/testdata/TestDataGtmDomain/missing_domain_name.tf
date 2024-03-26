provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_domain" "domain" {
}