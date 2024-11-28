provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_rules_builder" "default" {
  rules_v2024_01_09 {
    name          = "default"
    is_secure     = false
    comments      = "test"
    uuid          = "test"
    template_uuid = "test"
    template_link = "test"

    behavior {
      origin {
        cache_key_hostname = "ORIGIN_HOSTNAME"
        compress           = true
        custom_certificate_authorities {}
        custom_certificates {}
        enable_true_client_ip            = true
        forward_host_header              = "REQUEST_HOST_HEADER"
        http_port                        = 80
        https_port                       = 443
        origin_certs_to_honor            = "COMBO"
        origin_sni                       = true
        origin_type                      = "CUSTOMER"
        standard_certificate_authorities = []
        true_client_ip_client_setting    = false
        true_client_ip_header            = "True-Client-IP"
        use_unique_cache_key             = false
        verification_mode                = "CUSTOM"
      }
    }
  }
}
