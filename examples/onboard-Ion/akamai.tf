provider "akamai" {
  #Credentials can be provided inline using services such as Terraform Vault or by via an .edgerc file
  edgerc       = "../../.edgerc"
  papi_section = "papi"
  #  property {
  #        host = "${var.akamai_host}"
  #        access_token = "${var.akamai_access_token}"
  #        client_token = "${var.akamai_client_token}"
  #        client_secret = "${var.akamai_client_secret}"
  #    }
}

data "akamai_group" "group" {
  name = "grp_xxxx"
}

data "akamai_contract" "contract" {
  group = "${data.akamai_group.group.name}"
}

data "akamai_cp_code" "cp_code" {
  name     = "xxxxxx"
  group    = data.akamai_group.group.id
  contract = data.akamai_contract.contract.id
}

data "template_file" "rule_template" {
  template = "${file("${path.module}/rules/rules.json")}"
  vars = {
    snippets = "${path.module}/rules/snippets"
  }
}

data "template_file" "rules" {
  template = "${data.template_file.rule_template.rendered}"
  vars = {
    tdenabled = var.tdenabled
  }
}

resource "akamai_edge_hostname" "example-property" {
  product       = "prd_xxxx"
  contract      = data.akamai_contract.contract.id
  group         = data.akamai_group.group.id
  edge_hostname = "xxxx.edgesuite.net"
}

resource "akamai_property" "example-property" {
  name        = "example.mydomain.com"
  cp_code     = data.akamai_cp_code.cp_code.id
  contact     = ["me@mydomain.com"]
  contract    = data.akamai_contract.contract.id
  group       = data.akamai_group.group.id
  product     = "prd_xxxx"
  rule_format = "latest"
  hostnames = {
    "example.mydomain.com" = "${akamai_edge_hostname.example-property.edge_hostname}",
  }
  rules     = "${data.template_file.rules.rendered}"
  is_secure = true
}

resource "akamai_property_activation" "example-property" {
  property = "${akamai_property.example-property.id}"
  contact  = ["me@mydomain.com"]
  network  = "STAGING"
}
