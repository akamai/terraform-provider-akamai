provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  group_id    = "grp_0"
  contract_id = "ctr_0"
  product_id  = "prd_0"
  # Fetch the newly created property
  depends_on = [
    akamai_property.test
  ]
  rules = jsonencode(
    {
      "rules" : {
        "behaviors" : [
          {
            "name" : "origin",
            "options" : {
              "hostname" : "1.2.3.4",
              "httpPort" : 80,
              "httpsPort" : 443
            }
          }
        ]
      }
  })
}
