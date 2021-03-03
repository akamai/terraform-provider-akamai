provider "akamai" {
  edgerc = "~/.edgerc_network"
}


resource "akamai_networklist_subscription" "test" {
   recipients = ["test@email.com"]
    network_list = ["79536_MARTINNETWORKLIST"] 
   }

