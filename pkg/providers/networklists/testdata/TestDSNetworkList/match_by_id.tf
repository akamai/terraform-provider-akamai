provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_networklist_network_lists" "test" {
  name = "40996_ARTYLABWHITELIST"
  type = "IP"
}

