resource "akamai_cloudwrapper_configuration" "test" {
  config_name  = "testname"
  contract_id  = "ctr_123"
  property_ids = ["200200200"]
  comments     = "test"
  location {
    traffic_type_id = 1
    comments        = "test"
    capacity {
      value = 1
      unit  = "GB"
    }
  }
}