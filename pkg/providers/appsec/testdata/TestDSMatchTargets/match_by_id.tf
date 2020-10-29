provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_match_targets" "test" {
    config_id = 43253
    version = 7
    type =  "website"
    is_negative_path_match =  false
    is_negative_file_extension_match =  true
    default_file = "BASE_MATCH"
    hostnames =  ["example.com","www.example.net","m.example.com"]
    //file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
    //file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
    security_policy = "f1rQ_106946"
 
    bypass_network_lists = ["888518_ACDDCKERS","1304427_AAXXBBLIST"]
    
}
