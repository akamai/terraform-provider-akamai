provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_security_policy_clone" "test" {
    config_id = 43253
    version = 15 
    
    create_from_security_policy = "LNPD_76189"
    policy_name = "Cloned Test for Launchpad 15"
    policy_prefix = "LN" 
   }

