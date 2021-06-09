provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_eval_protect_host" "test" { 
  config_id = 43253
  hostnames = ["example.com"]
}


