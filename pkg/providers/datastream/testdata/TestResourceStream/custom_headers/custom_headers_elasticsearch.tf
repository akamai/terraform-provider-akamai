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

  elasticsearch_connector {
    connector_name      = "elasticsearch_connector_name"
    endpoint            = "endpoint"
    index_name          = "index_name"
    user_name           = "user_name"
    password            = "password"
    content_type        = "content_type"
    custom_header_name  = "custom_header_name"
    custom_header_value = "custom_header_value"
  }
}
