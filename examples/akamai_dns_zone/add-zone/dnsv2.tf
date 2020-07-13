provider "akamai" {
  dns_section = "dns"
}

locals {
  pzone = "example_primary.net"
  szone = "example_secondary.net"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

// Example primary zone resource
// NOTE: Please review the Provider Getting Started documentation before creating a Primary zone
//
resource "akamai_dns_zone" "test_primary_zone" {
	contract = data.akamai_contract.contract.id
        group = data.akamai_group.group.id
	zone = local.pzone
	type = "PRIMARY"
	comment =  "This is a test  primary zone"
	sign_and_serve = false
}

resource "akamai_dns_zone" "test_secondary_zone" {
        contract = data.akamai_contract.contract.id
        group = data.akamai_group.group.id
        zone = local.szone
        masters = ["1.2.3.4" , "1.2.3.5"]
        type = "secondary"
        comment =  "This is a test secondary zone"
        sign_and_serve = false
}

data "akamai_authorities_set" "ns" {
  contract = data.akamai_contract.contract.id
}


output "authorities" {
  value = join(",", data.akamai_authorities_set.ns.authorities)
}


