variable "AKAMAI_HOST" {
  type = string
}

variable "AKAMAI_ACCESS_TOKEN" {
  type = string
}

variable "AKAMAI_CLIENT_TOKEN" {
  type = string
}

variable "AKAMAI_CLIENT_SECRET" {
  type = string
}

provider "akamai" {
  property {
    host = "${var.AKAMAI_HOST}"
    access_token = "${var.AKAMAI_ACCESS_TOKEN}"
    client_token = "${var.AKAMAI_CLIENT_TOKEN}"
    client_secret = "${var.AKAMAI_CLIENT_SECRET}"
  }
}

resource "akamai_property" "terraform_example" {
  name = "terraform_example"
  contact = [
    "dshafik@akamai.com"
  ]
  product = "prd_SPM"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"
  cp_code = "cpc_846642"

  hostnames = {
    "terraform.example.org" = "${akamai_edge_hostname.ehn.edge_hostname}"
    "terraform.example.com" = "${akamai_edge_hostname.ehn.edge_hostname}"
  }

  variables = "${akamai_property_variables.origin.json}"
  rules = "${data.local_file.rules.content}"
}

resource "akamai_edge_hostname" "ehn" {
  edge_hostname = "terraform.example.org.edgesuite.net"

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
      value       = "origin.example.org"
      description = "Terraform Demo Origin"
      hidden      = true
      sensitive   = false
    }
  }
}