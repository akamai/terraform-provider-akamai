provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_domains" "domains" {
}