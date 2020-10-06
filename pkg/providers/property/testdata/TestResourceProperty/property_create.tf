
provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "prop" {
  name = "property_name"
  contact = ["test@akamai.com"]
  product = "prd_2"
  cp_code = "cpc_1"
  contract = "ctr_2"
  group = "grp_2"

hostnames = {
  "cnamefrom" = "akamai.edgesuite.net"
}

rule_format = "rule_format"

}
