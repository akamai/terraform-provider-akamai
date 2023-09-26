provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudwrapper_configuration" "test" {
  config_name         = "testname"
  contract_id         = "ctr_234"
  property_ids        = ["200200200"]
  notification_emails = ["test@akamai.com"]
  comments            = "test"
  retain_idle_objects = false
  location {
    traffic_type_id = 1
    comments        = "test"
    capacity {
      value = 1
      unit  = "GB"
    }
  }
}