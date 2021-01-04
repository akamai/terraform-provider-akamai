provider "akamai" {
  edgerc = "~/.edgerc_network"
}


resource "akamai_networklist_network_list_description" "test" {
   uniqueid =  "79536_MARTINNETWORKLIST"
     name = "Martin Network List"
     description =  "Notes about this network list"  
   }

