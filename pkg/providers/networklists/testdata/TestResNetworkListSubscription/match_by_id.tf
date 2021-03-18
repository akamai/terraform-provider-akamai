provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_networklist_subscription" "test" {
   recipients = ["test@email.com"]
    network_list = ["79536_MARTINNETWORKLIST"] 
   }

