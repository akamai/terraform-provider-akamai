provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_lineages" "test" {
  contract_id        = "C-0N7RAC7"
  lineage_name       = "example"
  lineage_ids        = [500001, 500002]
  secure_network     = "ENHANCED_TLS"
  stack_mode         = "MULTIPLE_STACK"
  lineage_type       = "MULTIPLE_GENERATION"
  domain             = "example.com"
  generation_status  = ["READY_FOR_USE", "ACTIVE"]
  expiring_in_days   = 30
  key_type           = "RSA"
  issuer             = "Test Certificate Authority"
  expand_generations = true
  sort               = "-modifiedDate"
}
