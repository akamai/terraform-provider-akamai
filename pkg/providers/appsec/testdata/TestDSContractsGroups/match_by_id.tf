provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_contracts_groups" "test" {
}
