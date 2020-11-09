provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_penalty_box" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
    action = "alert" 
    penalty_box_protection = true
}

output "appsecpenaltybox" {
  value = akamai_appsec_penalty_box.test.output_text
}
