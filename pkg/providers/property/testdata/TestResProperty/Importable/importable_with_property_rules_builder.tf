provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_builder" "default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    behavior {
      m_pulse {
        config_override = <<EOT
no new line
%{~if false}trim redundant new line%{endif~}
EOT
      }
    }
    behavior {
      m_pulse {
        config_override = <<EOT

%{~if false}trim redundant new line%{endif~}
EOT
      }
    }
    behavior {
      m_pulse {
        config_override = <<EOT

	line with new line before and after + tab

%{~if false}trim redundant new line%{endif~}
EOT
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
