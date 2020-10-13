provider "akamai" {
  edgerc = "~/.edgerc"
}

locals {
  contract = "9-CONTRACT"
  group    = "12345"
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

/*
resource "akamai_gtm_domain" "testdomain" {
        name = "gtm_terra_testdomain.akadns.net"
        type = "weighted"
        contract = "1-2ABCDEF"
        comment =  "Test"
        group     = "123ABC" 
        load_imbalance_percentage = 10
}
*/

resource "akamai_gtm_resource" "tfexample_resource_1" {
  domain           = local.gtmTestDomain 
  name             = "tfexample_resource_1"
  aggregation_type = "latest"
  type             = "XML load object via HTTP"
  resource_instance {
    datacenter_id           = 3131 
    use_default_load_object = false
    load_object             = "/test1"
    load_servers            = ["1.2.3.4"]
    load_object_port        = 80
  }
  wait_on_complete = false
}
