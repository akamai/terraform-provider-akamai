provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edgekv" "test" {
  namespace_name       = "DevExpTest"
  network              = "staging"
  group_id             = 4321
  retention_in_seconds = 86401
}