provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_authorities_set" "test" {
  contract = "ctr_xxxTestxxx"
}

output "authorities" {
  value = join(",", data.akamai_authorities_set.test.authorities)
}