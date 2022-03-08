provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_properties_search" "test" {
  key   = "hostname"
  value = "www.example.com"
}
