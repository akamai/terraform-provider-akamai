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
    upload_file_prefix = "prefix_updated"
    upload_file_suffix = "suf_updated"
  }

  contract_id = "test_contract"
  dataset_fields_ids = [
    2000, 1002, 2001, 1001
  ]
  email_ids = [
    "test_email1_updated@akamai.com",
    "test_email2@akamai.com",
  ]
  group_id = 1337
  property_ids = [
    1,
    2,
    3
  ]
  stream_name   = "test_stream_with_updated"
  stream_type   = "RAW_LOGS"
  template_name = "EDGE_LOGS"

  s3_connector {
    access_key        = "s3_test_access_key"
    bucket            = "s3_test_bucket_updated"
    connector_name    = "s3_test_connector_name_updated"
    path              = "s3_test_path"
    region            = "s3_test_region"
    secret_access_key = "s3_test_secret_key"
  }
}
