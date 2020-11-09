provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_match_targets" "test" {
    config_id = 43253
    version = 7
    
    
}
