provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_dns_record" "a_record" {
	zone = "exampleterraform.io"
	name = "exampleterraform.io"
	recordtype =  "A"
	active = true
	ttl = 300
	target = ["10.0.0.4","10.0.0.5"]
}

resource "akamai_dns_record" "soa_record" {
	zone = "exampleterraform.io"
	name = "@"
	recordtype =  "SOA"
	active = true
	ttl = 300
	name_server = "ns1.exampleterraform.io"
	email_address = "root+update@exampleterraform.io"
	refresh = 3600
	retry = 600
	expiry = 3600
	nxdomain_ttl = 3600
}