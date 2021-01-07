provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_reputation_analysis" "test" {
  config_id = 43253
  version = 12
  security_policy_id = "AAAA_81230"
 forward_to_http_header =  true
  forward_shared_ip_to_http_header_siem = true
}

