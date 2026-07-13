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

  netstorage_connector {
    display_name      = "display_name"
    secret_access_key = "secret_access_key"
    path              = "path"
    user_name         = "user_name"
    domain_prefix     = "domain_prefix"
  }
}
