provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules_builder" "default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    custom_override {
      name        = "test"
      override_id = "test"
    }
    advanced_override     = "test"
    comments              = "test"
    criteria_must_satisfy = "test"
    uuid                  = "test"
    template_uuid         = "test"
    template_link         = "test"
    criteria_locked       = true

    behavior {
      caching {
        behavior = "NO_STORE"
      }
    }

    children = [
      data.akamai_property_rules_builder.content_compression.json,
    ]
  }
}

data "akamai_property_rules_builder" "content_compression" {
  rules_v2023_01_05 {
    name                  = "Content Compression"
    criteria_must_satisfy = "all"
    variable {
      name        = "PMUSER_TESTSTR"
      description = "DSTR"
      value       = "STR"
      hidden      = false
      sensitive   = true
    }
    behavior {
      gzip_response {
        behavior = "ALWAYS"
      }
    }
  }
}
