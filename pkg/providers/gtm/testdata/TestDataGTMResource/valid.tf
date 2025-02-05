provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_resource" "my_gtm_resource" {
  domain        = "gtm_terra_testdomain.akadns.net"
  resource_name = "tfexample_resource_1"
}