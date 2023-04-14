provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_builder" "default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    behavior {
      m_pulse {
        config_override = trimsuffix(<<EOT
no new line
EOT
        , "\n")
      }
    }
    behavior {
      m_pulse {
        config_override = trimsuffix(<<EOT

EOT
        , "\n")
      }
    }
    behavior {
      m_pulse {
        config_override = trimsuffix(<<EOT

	line with new line before and after + tab

EOT
        , "\n")
      }
    }
  }
}

resource "akamai_property" "test" {
  name        = "test_property"
  group_id    = "grp_0"
  contract_id = "ctr_0"
  product_id  = "prd_0"

  hostnames {
    cname_to               = "to.test.domain"
    cname_from             = "from.test.domain"
    cert_provisioning_type = "DEFAULT"
  }
  rules = data.akamai_property_rules_builder.default.json
}
