provider "akamai" {
  edgerc = "~/.edgerc"
}

locals {
  gtmTestDomain = "gtm_terra_testdomain.akadns.net"
}

resource "akamai_gtm_resource" "tfexample_resource_1" {
  domain           = local.gtmTestDomain
  name             = "tfexample_resource_1"
  aggregation_type = "latest"
  type             = "XML load object via HTTP"
  resource_instance {
    datacenter_id           = 3132
    use_default_load_object = false
    load_object             = "/test2"
    load_servers            = ["1.2.3.5"]
    load_object_port        = 80
  }
  wait_on_complete = true
}
