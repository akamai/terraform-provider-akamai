provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  contract_id = "ctr_1"
  group_id    = "grp_2"
  product_id  = "prd_3"

  rules       = data.akamai_property_rules_builder.rules.json
  rule_format = data.akamai_property_rules_builder.rules.rule_format

  hostnames {
    cert_provisioning_type = "DEFAULT"
    cname_from             = "from.test.domain"
    cname_to               = "to.test.domain"
    cname_type             = "EDGE_HOSTNAME"
  }
}

data "akamai_property_rules_builder" "rules" {
  rules_v2023_01_05 {
    name = "default"

    variable {
      name        = "PMUSER_AKHOST"
      description = "Original Host Header"
      value       = ""
      hidden      = false
      sensitive   = false
    }
    variable {
      name        = "PMUSER_ENV"
      description = "environment indicator"
      value       = "DEV"
      hidden      = false
      sensitive   = false
    }
    variable {
      name        = "PMUSER_PATH"
      description = "Original Request Path"
      value       = ""
      hidden      = false
      sensitive   = false
    }
    variable {
      name        = "PMUSER_GRN"
      description = "Global Request Number"
      value       = ""
      hidden      = false
      sensitive   = false
    }
    variable {
      name        = "PMUSER_ACLBLOCKED"
      description = "User is blocked by access control"
      value       = "false"
      hidden      = false
      sensitive   = false
    }

    behavior {
      origin {
        http_port = 80
      }
    }
  }
}
