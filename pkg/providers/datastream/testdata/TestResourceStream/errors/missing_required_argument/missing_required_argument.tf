provider "akamai" {
    edgerc = "~/.edgerc"
}

resource "akamai_datastream" "s" {
    active = true
    config {
        delimiter          = "SPACE"
        format             = "STRUCTURED"
        frequency {
            time_in_sec = 30
        }
        upload_file_prefix = "pre"
        upload_file_suffix = "suf"
    }

    contract_id        = "test_contract"
    dataset_fields_ids = [
        1000,
        1001,
        1002
    ]
    email_ids          = [
        "test_email1@akamai.com",
        "test_email2@akamai.com",
    ]
    group_id           = 1337
    property_ids       = [
        1,
        2,
        3
    ]
    stream_type        = "RAW_LOGS"
    template_name      = "EDGE_LOGS"

    s3_connector {
        access_key        = "s3_test_access_key"
        bucket            = "s3_test_bucket"
        connector_name    = "s3_test_connector_name"
        path              = "s3_test_path"
        region            = "s3_test_region"
        secret_access_key = "s3_test_secret_key"
    }
}
