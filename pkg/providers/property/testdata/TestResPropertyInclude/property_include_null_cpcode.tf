provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_include" "test" {
  contract_id = "ctr_123"
  group_id    = "grp_123"
  product_id  = "prd_test"
  name        = "test include"
  type        = "MICROSERVICES"
  rule_format = data.akamai_property_rules_builder.rules_with_null.rule_format
  rules       = data.akamai_property_rules_builder.rules_with_null.json
}

data "akamai_property_rules_builder" "rules_with_null" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    behavior {
      cp_code {
        value {
          description = "CliTerraformCPCode"
          id          = 1047836
          name        = "DevExpCliTerraformPapiTest"
          products    = ["Web_App_Accel", ]
        }
      }
    }
  }
}