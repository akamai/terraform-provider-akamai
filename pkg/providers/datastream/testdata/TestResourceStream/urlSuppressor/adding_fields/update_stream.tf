provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_datastream" "s" {
  active = false
  delivery_configuration {
    field_delimiter = "SPACE"
    format          = "STRUCTURED"
    frequency {
      interval_in_secs = 30
    }
  }

  contract_id = "test_contract"
  dataset_fields = [
    1001
  ]
  group_id = 1337
  properties = [
    1,
  ]
  stream_name = "test_stream"

  sumologic_connector {
    collector_code      = "collector_code"
    display_name        = "display_name"
    endpoint            = "endpoint"
    custom_header_name  = "custom_header_name"
    custom_header_value = "custom_header_value"
  }
}
