provider "akamai" {
    edgerc = "../../test/edgerc"
}

data "akamai_gtm_domain" "domain" {
}