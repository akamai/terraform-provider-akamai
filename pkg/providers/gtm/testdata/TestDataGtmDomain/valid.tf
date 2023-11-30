provider "akamai" {
     edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_domain" "domain" {
    name = "test.cli.devexp-terraform.akadns.net"
}