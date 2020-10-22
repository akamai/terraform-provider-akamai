provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_authorities_set" "test" {
	
}
  
output "authorities" {
	value = "${join(",", data.akamai_authorities_set.test.authorities)}"
}