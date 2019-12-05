provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "global"
}

resource "akamai_property" "terraform_example" {
  name = "terraform_example1"
  contact = [
    "martin@akava.io"
  ]
  product = "prd_SPM"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"
  cp_code = "cpc_846642"

  hostnames = {
    "terraform.example1.org" = "${akamai_edge_hostname.ehn.edge_hostname}"
    "terraform.example1.com" = "${akamai_edge_hostname.ehn.edge_hostname}"
  }

  rule_format = "v2019-07-25"
  variables = "${akamai_property_variables.origin.json}"
  rules = "${data.local_file.rules.content}"
}

resource "akamai_edge_hostname" "ehn" {
  edge_hostname = "terraform.example1.org.edgesuite.net"

  product = "prd_SPM"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"

  ipv4 = true
  ipv6 = true
}

data "akamai_contract" "contract" {
  group = "${data.akamai_group.group.name}"
}

data "akamai_group" "group" {
  name = "Terraform Provider"
}

data "local_file" "rules" {
  filename = "rules.json"
}

resource "akamai_property_variables" "origin" {
  variables {
    variable {
      name        = "PMUSER_ORIGIN"
      value       = "origin.example1.org"
      description = "Terraform Demo Origin"
      hidden      = true
      sensitive   = false
    }
  }
}