provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_datastream" "s" {
  active = false
  config {
    format = "JSON"
    frequency {
      time_in_sec = 30
    }
  }

  contract_id = "test_contract"
  dataset_fields_ids = [
    1001
  ]
  group_id = 1337
  property_ids = [
    1,
  ]
  stream_name   = "test_stream"
  stream_type   = "RAW_LOGS"
  template_name = "EDGE_LOGS"

  new_relic_connector {
    connector_name      = "new_relic_connector_name"
    endpoint            = "endpoint"
    auth_token          = "auth_token"
    content_type        = "content_type"
    custom_header_name  = "custom_header_name"
    custom_header_value = "custom_header_value"
  }
}
