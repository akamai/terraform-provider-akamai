provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  contract_id = "ctr_1"
  group_id    = "grp_2"
  product_id  = "prd_3"
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
