provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_networklist_description" "test" {
  network_list_id = "2275_VOYAGERCALLCENTERWHITELI"
  name            = "Voyager Call Center Whitelist"
  description     = "Notes about this network list"
}

