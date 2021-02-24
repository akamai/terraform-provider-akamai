provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_bypass_network_lists" "test" {
config_id = 43253
    version = 7

}



