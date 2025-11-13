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

  s3_compatible_connector {
    bucket            = "bucket"
    display_name      = "display_name"
    access_key        = "access_key"
    region            = "region"
    secret_access_key = "secret_access_key"



    path = "path"
  }
}
