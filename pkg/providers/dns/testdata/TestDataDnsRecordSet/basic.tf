provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_dns_record_set" "test" {
  zone        = "exampleterraform.io"
  host        = "exampleterraform.io"
  record_type = "A"
}
