provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_datastream" "s" {
  active = false
  delivery_configuration {
    format = "JSON"
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

  https_connector {
    authentication_type = "BASIC"
    display_name        = "HTTPS connector name"
    compress_logs       = true
    content_type        = "content_type"
    custom_header_name  = "custom_header_name"
    custom_header_value = "custom_header_value"
    endpoint            = "https_connector_url"
    user_name           = "username"
    password            = "password"
  }
}
