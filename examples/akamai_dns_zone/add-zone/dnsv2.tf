provider "akamai" {
  dns_section = "dns"
}

locals {
  zone = "akavadev.net"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_dns_zone" "test_zone" {
	contract = data.akamai_contract.contract.id
	zone = local.zone
	#masters = ["1.2.3.4" , "1.2.3.5"]
	type = "PRIMARY"
	comment =  "This is a test zone"
	group     = data.akamai_group.group.id
	sign_and_serve = false
}

data "akamai_authorities_set" "ns" {
  contract = data.akamai_contract.contract.id
}


output "authorities" {
  value = join(",", data.akamai_authorities_set.ns.authorities)
}


