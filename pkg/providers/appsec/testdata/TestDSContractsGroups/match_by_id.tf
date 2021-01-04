provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_contract_groups" "contract_groups" {
}
