provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_dns_record_set" "test" {
  zone        = "exampleterraform.io"
  host        = "exampleterraform.io"
  record_type = "A"
}

output "test_addrs" {
  value = join(",", data.akamai_dns_record_set.test.rdata)
}