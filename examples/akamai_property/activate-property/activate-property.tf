provider "akamai" {
     edgerc = "~/.edgerc"
     papi_section = "papi_section"
}

variable "activate" {
	default = true
}

resource "akamai_property_activation" "dshafik_sandbox" {

        name = "akavadeveloper.com"
        contact = ["martin@akava.io"]
        hostname = ["akavadeveloper.com"]
        contract = "${data.akamai_contract.our_contract.id}"
        group =  "${data.akamai_group.our_group.id}"
        network = "STAGING"
        activate = "${var.activate}"
}

data "akamai_group" "our_group" {
    name = "Davey Shafik"
}

output "groupid" {
  value = "${data.akamai_group.our_group.id}"
}


data "akamai_contract" "our_contract" {
    name = "Davey Shafik"
}

output "contractid" {
  value = "${data.akamai_contract.our_contract.id}"
}
