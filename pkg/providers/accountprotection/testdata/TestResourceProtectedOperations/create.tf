provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_apr_protected_operations" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  operation_id       = "b85e3eaa-d334-466d-857e-33308ce416be"

  protected_operation = jsonencode(
    {
      "testKey" : "testValue"
    }
  )
}
