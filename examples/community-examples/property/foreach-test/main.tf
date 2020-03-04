provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "papi"
}


data "akamai_group" "group" {
  name = "Ian Cass"
}

data "akamai_contract" "contract" {
  group = data.akamai_group.group.name
}

data "template_file" "rule_template" {
	template = file("${path.module}/rules/rules.json")
	vars = {
		snippets = "${path.module}/rules/snippets"
	}
}

data "template_file" "rules" {
	for_each = var.customers

	template = data.template_file.rule_template.rendered
	vars = {
		username = each.value.username
		password = each.value.password
	}
}

resource "akamai_cp_code" "cpcode" {

    for_each = var.customers

    product  = "prd_Site_Accel"
    contract = data.akamai_contract.contract.id
    group = data.akamai_group.group.id
    name	= each.key
}

resource "akamai_edge_hostname" "edge_hostname" {

    product  = "prd_Site_Accel"
    contract = data.akamai_contract.contract.id
    group = data.akamai_group.group.id
    edge_hostname = "test.wheep.co.uk.edgesuite.net"
}

resource "akamai_property" "property" {
  
  for_each = var.customers 

  name        = each.key
  cp_code     = akamai_cp_code.cpcode[each.key].id
  contact     = [""]
  contract = data.akamai_contract.contract.id
  group = data.akamai_group.group.id
  product     = "prd_Site_Accel"
  rule_format = "v2018-02-27"

  hostnames    = {
	"${each.key}" = akamai_edge_hostname.edge_hostname.edge_hostname
  }
  rules       = data.template_file.rules[each.key].rendered
  is_secure = true

}

resource "akamai_property_activation" "activation" {
  for_each = var.customers 

  property = akamai_property.property[each.key].id
  contact = ["icass@akamai.com"]
  network = upper(var.env)
  activate = true
}
