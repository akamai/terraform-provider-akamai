provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_datastream" "s" {
  active = true
  delivery_configuration {
    field_delimiter = "SPACE"
    format          = "STRUCTURED"
    frequency {
      interval_in_secs = 30
    }
    upload_file_prefix = "pre"
    upload_file_suffix = "suf"
  }

  contract_id    = "test_contract"
  dataset_fields = [1001]
  group_id       = 1337
  properties     = [1]
  stream_name    = "test_stream"

  oracle_connector {
    access_key        = "access_key"
    bucket            = "bucket"
    display_name      = "display_name"
    namespace         = "namespace"
    path              = "path"
    region            = "region"
    secret_access_key = "secret_access_key"
  }
}
