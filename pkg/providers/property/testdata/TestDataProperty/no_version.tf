provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_property" "prop" {
  name = "property_name"
}
