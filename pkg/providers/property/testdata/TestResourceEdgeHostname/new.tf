provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract      = "ctr_2"
  group         = "grp_2"
  product       = "prd_2"
  edge_hostname = "test.aka"
  ip_behavior   = "IPV4"
  use_cases     = <<-EOF
  [
  {
    "option": "BACKGROUND",
    "type": "GLOBAL",
    "useCase": "Download_Mode"
  }
  ]
  EOF
}

output "edge_hostname" {
  value = akamai_edge_hostname.edgehostname.edge_hostname
}
