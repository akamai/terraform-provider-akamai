provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
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

  elasticsearch_connector {
    display_name = "elasticsearch_connector_name"
    endpoint     = "endpoint"
    index_name   = "index_name"
    user_name    = "user_name"
    password     = "password"
    content_type = "content_type"
    tls_hostname = "tls_hostname"
    ca_cert      = "ca_cert"
    client_cert  = "client_cert"
    client_key   = "client_key"
  }
}
