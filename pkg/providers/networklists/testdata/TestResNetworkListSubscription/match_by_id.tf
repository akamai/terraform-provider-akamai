provider "akamai" {
  edgerc = "~/.edgerc_network"
}


resource "akamai_networklist_network_list_subscription" "test" {
   recipients = ["test@email.com"]
    network_list_ids = ["79536_MARTINNETWORKLIST"] 
   }

