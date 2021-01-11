provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_eval_hostnames" "test" {
  config_id = 43253
    version = 7
}

