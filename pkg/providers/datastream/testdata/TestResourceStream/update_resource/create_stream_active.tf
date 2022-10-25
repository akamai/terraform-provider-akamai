provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_datastream" "s" {
  active = true
  config {
    delimiter = "SPACE"
    format    = "STRUCTURED"
    frequency {
      time_in_sec = 30
    }
    upload_file_prefix = "pre"
    upload_file_suffix = "suf"
  }

  contract_id        = "test_contract"
  dataset_fields_ids = [1001]
  group_id           = 1337
  property_ids       = [1]
  stream_name        = "test_stream"
  stream_type        = "RAW_LOGS"
  template_name      = "EDGE_LOGS"

  oracle_connector {
    access_key        = "access_key"
    bucket            = "bucket"
    connector_name    = "connector_name"
    namespace         = "namespace"
    path              = "path"
    region            = "region"
    secret_access_key = "secret_access_key"
  }
}
