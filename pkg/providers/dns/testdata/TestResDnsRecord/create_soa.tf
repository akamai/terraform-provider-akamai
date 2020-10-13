provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_dns_record" "soa_record" {
	zone = "exampleterraform.io"
	name = "@"
	recordtype =  "SOA"
	active = true
	ttl = 300
	name_server = "ns1.exampleterraform.io"
	email_address = "root@exampleterraform.io"
	refresh = 3600
	retry = 600
	expiry = 3600
	nxdomain_ttl = 3600
}