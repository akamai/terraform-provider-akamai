provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_custom_rules" "test" {
  config_id = 43253

}
output "appseccustomrules" {
  value = data.akamai_appsec_custom_rules.test.output_text
}
