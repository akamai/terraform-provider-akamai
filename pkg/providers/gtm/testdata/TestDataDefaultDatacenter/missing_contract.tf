provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_gtm_default_datacenter" "test" {
	
}
  
output "datacenter_id" {
	value = "${join(",", data.akamai_gtm_default_datacenter.test.datacenter_id)}"
}
