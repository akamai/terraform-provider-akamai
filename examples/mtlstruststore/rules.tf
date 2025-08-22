# This example presents a sample workflow for building property rules with the mTLS Truststore behavior enabled.
#
# Before applying this example, make changes to the attribute values according to your needs.
# To configure base settings for mTLS Truststore, use the `enforce_mtls_settings` behavior and the `request_header` criterion.
# For more custom settings, add additional features to your rules, including the `client_certificate` criterion, and the `client_certificate_auth` and `log_custom` behaviors.
#
# A successful operation creates property rules with the mTLS Truststore behavior enabled.
#
# You can use the created rules to enforce mTLS Truststore settings on the property.

data "akamai_property_rules_builder" "full_mtls_workflow_rule_default" {
  rules_v2024_02_12 {
    name      = "default"
    is_secure = false
    behavior {
      origin_characteristics {
        authentication_method       = "AUTOMATIC"
        authentication_method_title = ""
        country                     = "UNKNOWN"
        origin_location_title       = ""
      }
    }
    behavior {
      origin {
        cache_key_hostname            = "ORIGIN_HOSTNAME"
        compress                      = true
        enable_true_client_ip         = true
        forward_host_header           = "REQUEST_HOST_HEADER"
        hostname                      = "www.example.com"
        http_port                     = 80
        https_port                    = 443
        ip_version                    = "IPV4"
        max_tls_version               = "DYNAMIC"
        min_tls_version               = "DYNAMIC"
        origin_certificate            = ""
        origin_sni                    = true
        origin_type                   = "CUSTOMER"
        ports                         = ""
        tls_version_title             = ""
        true_client_ip_client_setting = false
        true_client_ip_header         = "True-Client-IP"
        verification_mode             = "PLATFORM_SETTINGS"
      }
    }
    behavior {
      cp_code {
        value {
          id = akamai_cp_code.cp_code.id
        }
      }
    }
    behavior {
      cache_key_query_params {
        behavior = "IGNORE_ALL"
      }
    }
    behavior {
      http3 {
        enable = true
      }
    }
    behavior {
      caching {
        behavior = "NO_STORE"
      }
    }
    children = [
      data.akamai_property_rules_builder.full_mtls_workflow_rule_m_tls_settings_enforcement_base.json,
    ]
  }
}

data "akamai_property_rules_builder" "full_mtls_workflow_rule_m_tls_settings_enforcement_base" {
  rules_v2024_02_12 {
    name                  = "mTLS Settings Enforcement – Base"
    criteria_must_satisfy = "any"
    criterion {
      request_header {
        header_name                = "test_1"
        match_case_sensitive_value = true
        match_operator             = "IS_ONE_OF"
        match_wildcard_name        = false
        match_wildcard_value       = false
        values                     = ["ON", ]
      }
    }
    behavior {
      enforce_mtls_settings {
        certificate_authority_set = [akamai_mtlstruststore_ca_set_activation.ca_set_activation_production.ca_set_id, ]
        enable_auth_set           = true
        enable_deny_request       = false
        enable_ocsp_status        = false
      }
    }
    children = [
      data.akamai_property_rules_builder.full_mtls_workflow_rule_m_tls_settings_enforcement_custom.json,
    ]
  }
}

data "akamai_property_rules_builder" "full_mtls_workflow_rule_m_tls_settings_enforcement_custom" {
  rules_v2024_02_12 {
    name                  = "Client Certificate Settings – Custom"
    criteria_must_satisfy = "all"
    criterion {
      client_certificate {
        enforce_mtls           = true
        is_certificate_present = true
        is_certificate_valid   = "INVALID"
      }
    }
    behavior {
      client_certificate_auth {
        client_certificate_attributes               = []
        enable                                      = true
        enable_client_certificate_validation_status = true
        enable_complete_client_certificate          = true
      }
    }
    behavior {
      log_custom {
        log_custom_log_field = false
      }
    }
  }
}
