
resource "akamai_property" "akavadeveloper" {
    name = "akavadev1.com"
    contact = ["dshafik@akamai.com"]
    product = "prd_SPM"
    contract = "${data.akamai_contract.our_contract.id}" #"ctr_C-1FRYVV3"
    group = "${data.akamai_group.our_group.id}" #"grp_68817"
    cp_code = "${data.akamai_cp_codes.cp_codes.id}" #"cpc_846350"
    # Use one of hostnames or secure_hostnames, we may not need to separate
    #hostname = ["www.akavadev1.com"]
    #activate = false

    origin {
    is_secure = false
    hostname = "www.akava1.com"
    forward_hostname = "ORIGIN_HOSTNAME"
  }



    # rules takes a resulting JSON string with the full rules

    # Load from a data.local_file
    #rulesjson = "${jsonencode(data.local_file.akavadeveloper_json.content)}"

    # Generate from a resource.akamai_property_rules
    rules = "${akamai_property_rules.akavadeveloper.json}"
    #rules =   ["${jsonencode(local.rulesObject)}"]
}


resource "akamai_secure_edge_hostname" "akavadeveloper" {
    #name = "akavadev1.com.edgekey.net"
    product = "prd_SPM"
    contract = "${data.akamai_contract.our_contract.id}" #"ctr_C-1FRYVV3"
    group = "${data.akamai_group.our_group.id}" #"grp_68817"
    edge_hostname =  "www.akavadev1-1.edgekey.net"
    certenrollmentid = "12356666"
    ipv6 = true


}

/*
data "external" "generate_json" {
    program = ["/usr/bin/bash","${path.module}/merge.sh"]
   query = {
    p = "akavadev1.com"
   }
}
*/
/*
data "local_file" "akavadeveloper_json" {
    depends_on = ["data.external.generate_json"]
    #filename = "${path.module}/rules.json"
    filename =  "${path.module}/akavadev1.com/dist/akavadev1.com.papi.json"
}
*/

resource "akamai_property_rules" "akavadeveloper" {
    rules {
        behavior {
            name = "downstreamCache"
            option {
                key = "behavior"
                value = "TUNNEL_ORIGIN"
            }
        }


        rule {
            name = "Uncacheable Responses"
            comment = "Cache me outside"
            criteria {
                name = "cacheability"
                option {
                key = "matchOperator"
                value = "IS_NOT"
                }
                option {
                key = "value"
                value = "CACHEABLE"
                }
            }
            behavior {
                name = "downstreamCache"
                option {
                key = "behavior"
                value = "TUNNEL_ORIGIN"
                }
            }
            rule {
                name = "Uncacheable Responses"
                comment = "Child rule"
                criteria {
                    name = "cacheability"
                    option {
                        key = "matchOperator"
                        value = "IS_NOT"
                    }
                    option {
                        key = "value"
                        value = "CACHEABLE"
                    }
                }
                behavior {
                    name = "downstreamCache"
                    option {
                        key = "behavior"
                        value = "TUNNEL_ORIGIN"
                    }
                }
            }
        }
    }
}

resource "akamai_property_variable" "origin" {
    "name" = "ORIGIN"
    "value" = "origin.akavadeveloper.com"
    "description" = "Property Origin"
    "hidden" = false
    "sensitive" = false
    "fqname" = "origin.akavadev1.com"
}

/*
resource "akamai_cp_code" "akavadeveloper" {
    name = "PDM"
    product = "prd_SPM"
    contract = "ctr_C-1FRYVV3"
    group = "grp_130690" # "grp_68817"
   # cp_code = "cpc_846350"
}
*/

/*
resource "akamai_cp_code" "akavadeveloper" {
    name = "akavadev1.com"
    product = "prd_SPM"
    contract = "ctr_C-1FRYVV3"
    group = "grp_130690" # "grp_68817"
   # cp_code = "cpc_846350"
}
*/

data  "akamai_cp_codes" "cp_codes" {
    name = "akavadev1.com"
    contract = "ctr_C-1FRYVV3"
    group = "grp_130690" #grp_68817"

}

output "cp_codes" {
  value = "${data.akamai_cp_codes.cp_codes.id}"
}


data "akamai_group" "our_group" {
    name = "Terraform Provider" #"Davey Shafik"
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


output "json" {
  value = "${akamai_property_rules.akavadeveloper.json}"
}
