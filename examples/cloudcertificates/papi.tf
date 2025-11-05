# This example presents a sample workflow for configuring a property that uses a Cloud Certificate Manager (CCM) integration.
#
# Before applying this example, make changes to the attribute values according to your needs.
#
# A successful operation creates a property with a hostname bound to the cloud certificate, and activates that property on the `STAGING` and `PRODUCTION` environments.

data "akamai_property_rules_template" "rules" {
  template_file = abspath("${path.module}/property-snippets/main.json")
}

resource "akamai_property" "test" {
  name        = "test-property-name"
  contract_id = "C-0N7RAC7"
  group_id    = "grp_123"
  product_id  = "prd_12345"
  hostnames {
    cname_from             = "test.example.com"
    cname_to               = "test.example.com.edgekey.net"
    cert_provisioning_type = "CCM"
    ccm_certificates {
      rsa_cert_id = akamai_cloudcertificates_upload_signed_certificate.upload.certificate_id
    }
  }
  rule_format = "v2025-07-07"
  rules       = data.akamai_property_rules_template.rules.json
}

resource "akamai_property_activation" "test-staging" {
  property_id                    = akamai_property.test.id
  contact                        = ["test@example.com"]
  version                        = akamai_property.test.latest_version
  network                        = "STAGING"
  auto_acknowledge_rule_warnings = false
}

resource "akamai_property_activation" "test-production" {
  property_id                    = akamai_property.test.id
  contact                        = ["test@example.com"]
  version                        = akamai_property.test.latest_version
  network                        = "PRODUCTION"
  auto_acknowledge_rule_warnings = false
}