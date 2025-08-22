# This example presents a sample workflow for configuring a property that enforces the mTLS Truststore behavior.
#
# Before applying this example, make changes to the attribute values according to your needs.
#
# A successful operation creates an edge hostname, CP code, and property with the mTLS Truststore behavior enabled, and activates that property on `STAGING` and `PRODUCTION` environments.

resource "akamai_edge_hostname" "aka_edgehost" {
  contract_id   = "C-0N7RAC7"
  group_id      = "grp_123"
  product_id    = "prd_Site_Accel"
  edge_hostname = "www.example-hostname.edgesuite.net"
  ip_behavior   = "IPV4"
}

resource "akamai_cp_code" "cp_code" {
  contract_id = "C-0N7RAC7"
  group_id    = "grp_123"
  product_id  = "prd_Site_Accel"
  name        = "CP-Code-Name"
}

resource "akamai_property" "property" {
  name        = "Property-Name"
  product_id  = "prd_Site_Accel"
  contract_id = "C-0N7RAC7"
  group_id    = "grp_123"
  hostnames {
    cname_to               = akamai_edge_hostname.aka_edgehost.edge_hostname
    cname_from             = "www.example-hostname.edgesuite.net"
    cert_provisioning_type = "CPS_MANAGED"
  }
  rule_format = data.akamai_property_rules_builder.full_mtls_workflow_rule_default.rule_format
  rules       = data.akamai_property_rules_builder.full_mtls_workflow_rule_default.json
}

resource "akamai_property_activation" "property_activate_staging" {
  contact                        = ["jsmith@example.com"]
  network                        = "STAGING"
  property_id                    = akamai_property.property.id
  version                        = 1
  auto_acknowledge_rule_warnings = true
}

resource "akamai_property_activation" "property_activate_production" {
  contact                        = ["jsmith@example.com"]
  network                        = "PRODUCTION"
  property_id                    = akamai_property.property.id
  version                        = 1
  auto_acknowledge_rule_warnings = true
}