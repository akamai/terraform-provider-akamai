provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudaccess_key" "test" {
  access_key_name = "foo"
}